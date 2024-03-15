package u2i

import (
	"encoding/json"
	"infer-microservices/internal"
	faiss_index "infer-microservices/internal/faiss_gogofaster"
	"infer-microservices/internal/logs"
	"infer-microservices/internal/utils"
	"infer-microservices/pkg/ann/faiss"
	"infer-microservices/pkg/config_loader/faiss_config"
	"infer-microservices/pkg/config_loader/model_config"
	feature "infer-microservices/pkg/infer_features"
	"infer-microservices/pkg/infer_models/base_model"
	"net/http"
	"time"
)

//INFO: user to item DSSM model, used to calculate user & item similarity.

type Dssm struct {
	baseModel base_model.BaseModel // extend baseModel
	retNum    int
	modelType string
}

// retNum
func (d *Dssm) SetRetNum(retNum int) {
	d.retNum = retNum
}

func (d *Dssm) GetRetNum() int {
	return d.retNum
}

func (d *Dssm) SetBaseModel(baseModel base_model.BaseModel) {
	d.baseModel = baseModel
}

// modeltype
func (d *Dssm) SetModelType(modelType string) {
	d.modelType = modelType
}

func (d *Dssm) GetModelType() string {
	return d.modelType
}

func (d *Dssm) ModelInferSkywalking(model model_config.ModelConfig, requestId string, exposureList []string, r *http.Request, inferSample feature.ExampleFeatures, retNum int) (map[string][]map[string]interface{}, error) {
	response := make(map[string][]map[string]interface{}, 0)
	cacheKeyPrefix := requestId + d.baseModel.GetServiceConfig().GetServiceId() + d.baseModel.GetModelName() + d.baseModel.GetTensorName()
	// get rsp from cache.
	rspDataBytes, err := d.baseModel.GetBigCacheRsp().Get(cacheKeyPrefix)
	if err != nil {
		logs.Error(err)
	} else {
		err = json.Unmarshal(rspDataBytes, &response)
		if err == nil {
			return response, nil
		}
	}

	// get embedding from tfserving model.
	spanUnionEmFv, _, err := internal.GetTracer().CreateLocalSpan(r.Context())
	if err != nil {
		return nil, err
	}
	spanUnionEmFv.SetOperationName("get recall embedding func")
	spanUnionEmFv.Log(time.Now())

	embeddingVector, err := d.baseModel.Embedding(model, inferSample, d.baseModel.GetTensorName())
	if err != nil {
		logs.Error(requestId, time.Now(), err)
		return nil, err
	}
	spanUnionEmFv.Log(time.Now())
	spanUnionEmFv.End()
	logs.Debug(requestId, time.Now(), "embeddingVector:", embeddingVector)

	//Asynchronous RPC request, simultaneous processing of multiple recalls, reduce network cost
	mergeResult := make([]*faiss_index.ItemInfo, 0)
	faissIndexConfigs := d.baseModel.GetServiceConfig().GetFaissIndexConfigs()
	recallCh := make(chan []*faiss_index.ItemInfo, 100)
	for _, faissIndexConfig := range faissIndexConfigs.GetFaissIndexConfig() {
		go func(faissIndexConfig *faiss_config.FaissIndexConfig) {
			recallResult, err := faiss.FaissVectorSearch(faissIndexConfig, inferSample, *embeddingVector)
			if err != nil {
				logs.Error(err)
			}
			logs.Debug(requestId, time.Now(), "recall result:", recallResult)
			recallCh <- recallResult
		}(&faissIndexConfig)
	}

loop:
	for {
		select {
		case <-time.After(time.Millisecond * 100):
			break loop
		case recall := <-recallCh:
			for _, item := range recall {
				mergeResult = append(mergeResult, item)
			}
		}
	}
	close(recallCh)

	//format result.
	spanUnionEmOut, _, err := internal.GetTracer().CreateLocalSpan(r.Context())
	if err != nil {
		return nil, err
	}
	spanUnionEmOut.SetOperationName("get recall result func")

	spanUnionEmOut.Log(time.Now())
	recallRst, err := d.baseModel.InferResultFormat(&mergeResult, exposureList)
	if err != nil {
		return nil, err
	}
	spanUnionEmOut.Log(time.Now())
	spanUnionEmOut.End()

	if len(*recallRst) == 0 {
		logs.Error(requestId, time.Now(), "recall 0 item, check the faiss index plz. ")
		return nil, err
	}
	response["data"] = *recallRst
	logs.Debug(requestId, time.Now(), "format result:", mergeResult)

	d.baseModel.GetBigCacheRsp().Set(cacheKeyPrefix, []byte(utils.ConvertStructToJson(response)))

	return response, nil
}

func (d *Dssm) ModelInferNoSkywalking(model model_config.ModelConfig, requestId string, exposureList []string, inferSample feature.ExampleFeatures, retNum int) (map[string][]map[string]interface{}, error) {
	response := make(map[string][]map[string]interface{}, 0)
	cacheKeyPrefix := requestId + d.baseModel.GetServiceConfig().GetServiceId() + d.baseModel.GetModelName() + d.baseModel.GetTensorName()
	// get rsp from cache.
	rspDataBytes, err := d.baseModel.GetBigCacheRsp().Get(cacheKeyPrefix)
	if err != nil {
		logs.Error(err)
	} else {
		err = json.Unmarshal(rspDataBytes, &response)
		if err == nil {
			return response, nil
		}
	}
	// get embedding from tfserving model.
	embeddingVector, err := d.baseModel.Embedding(model, inferSample, d.baseModel.GetTensorName())
	if err != nil {
		return nil, err
	}
	logs.Debug(requestId, time.Now(), "embeddingVector:", embeddingVector)

	//Asynchronous RPC request, simultaneous processing of multiple recalls, reduce network cost
	mergeResult := make([]*faiss_index.ItemInfo, 0)
	faissIndexConfigs := d.baseModel.GetServiceConfig().GetFaissIndexConfigs()
	recallCh := make(chan []*faiss_index.ItemInfo, 100)
	for _, faissIndexConfig := range faissIndexConfigs.GetFaissIndexConfig() {
		go func(faissIndexConfig *faiss_config.FaissIndexConfig) {
			recallResult, err := faiss.FaissVectorSearch(faissIndexConfig, inferSample, *embeddingVector)
			if err != nil {
				logs.Error(err)
			}
			logs.Debug(requestId, time.Now(), "recall result:", recallResult)
			recallCh <- recallResult
		}(&faissIndexConfig)
	}

loop:
	for {
		select {
		case <-time.After(time.Millisecond * 100):
			break loop
		case recall := <-recallCh:
			for _, item := range recall {
				mergeResult = append(mergeResult, item)
			}
		}
	}
	close(recallCh)

	//format result.
	recallRst, err := d.baseModel.InferResultFormat(&mergeResult, exposureList)
	if err != nil {
		return nil, err
	}

	if len(*recallRst) == 0 {
		logs.Error("recall 0 item, check the faiss index plz. ")
		return nil, err
	}
	response["data"] = *recallRst
	logs.Debug(requestId, time.Now(), "format result:", mergeResult)

	d.baseModel.GetBigCacheRsp().Set(cacheKeyPrefix, []byte(utils.ConvertStructToJson(response)))

	return response, nil
}

func (d *Dssm) RecallSampleDirector(userId string, offlineFeature bool, onlineFeature bool, featureList []string) (feature.ExampleFeatures, error) {
	//panic("no not need to Implement interface methods")
	sample := feature.ExampleFeatures{}
	return sample, nil
}

func (d *Dssm) RankingSampleDirector(userId string, offlineFeature bool, onlineFeature bool, itemIdList []string, featureList []string) (feature.ExampleFeatures, error) {
	//panic("no not need to Implement interface methods")
	sample := feature.ExampleFeatures{}
	return sample, nil
}
