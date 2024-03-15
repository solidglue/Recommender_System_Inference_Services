package heavy_ranker

import (
	"encoding/json"
	"infer-microservices/internal"
	faiss_index "infer-microservices/internal/faiss_gogofaster"
	"infer-microservices/internal/logs"
	"infer-microservices/internal/utils"
	"infer-microservices/pkg/config_loader/model_config"
	feature "infer-microservices/pkg/infer_features"
	"infer-microservices/pkg/infer_models/base_model"
	"net/http"
	"time"
)

//INFO: short sequence model, such as hundreds videos a user played in the past days/hours.

type DIN struct {
	baseModel base_model.BaseModel // extend baseModel
	modelType string
}

func (d *DIN) SetBaseModel(baseModel base_model.BaseModel) {
	d.baseModel = baseModel
}

// modeltype
func (d *DIN) SetModelType(modelType string) {
	d.modelType = modelType
}

func (d *DIN) GetModelType() string {
	return d.modelType
}

func (d *DIN) ModelInferSkywalking(model model_config.ModelConfig, requestId string, exposureList []string, r *http.Request, inferSample feature.ExampleFeatures, retNum int) (map[string][]map[string]interface{}, error) {
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

	// get rank scores from tfserving model.
	rankResult := make([]*faiss_index.ItemInfo, 0)
	spanUnionEmFv, _, err := internal.GetTracer().CreateLocalSpan(r.Context())
	if err != nil {
		return nil, err
	}
	spanUnionEmFv.SetOperationName("get rank scores func")
	spanUnionEmFv.Log(time.Now())
	items, scores, err := d.baseModel.RankPredict(model, inferSample, d.baseModel.GetTensorName())
	if err != nil {
		return nil, err
	}
	spanUnionEmFv.Log(time.Now())
	spanUnionEmFv.End()
	logs.Debug(requestId, time.Now(), "unrank result:", inferSample)

	//build rank result whith tfserving.ItemInfo
	for idx := 0; idx < len(*items); idx++ {
		if idx > retNum {
			break
		}
		itemInfo := &faiss_index.ItemInfo{
			ItemId: (*items)[idx],
			Score:  (*scores)[idx],
		}
		rankResult = append(rankResult, itemInfo)
	}
	logs.Debug(requestId, time.Now(), "rank result:", inferSample)

	//format result.
	spanUnionEmOut, _, err := internal.GetTracer().CreateLocalSpan(r.Context())
	if err != nil {
		return nil, err
	}
	spanUnionEmOut.SetOperationName("get rank result func")
	spanUnionEmOut.Log(time.Now())
	rankRst, err := d.baseModel.InferResultFormat(&rankResult, exposureList)
	if err != nil {
		return nil, err
	}
	spanUnionEmOut.Log(time.Now())
	spanUnionEmOut.End()

	response["data"] = *rankRst
	logs.Debug(requestId, time.Now(), "format result", response)

	d.baseModel.GetBigCacheRsp().Set(cacheKeyPrefix, []byte(utils.ConvertStructToJson(response)))

	return response, nil
}

func (d *DIN) ModelInferNoSkywalking(model model_config.ModelConfig, requestId string, exposureList []string, inferSample feature.ExampleFeatures, retNum int) (map[string][]map[string]interface{}, error) {
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

	// get rank scores from tfserving model.
	rankResult := make([]*faiss_index.ItemInfo, 0)
	items, scores, err := d.baseModel.RankPredict(model, inferSample, d.baseModel.GetTensorName())
	if err != nil {
		logs.Error(requestId, time.Now(), err)
		return nil, err
	}
	logs.Debug(requestId, time.Now(), "unrank result:", inferSample)

	//build rank result whith tfserving.ItemInfo
	for idx := 0; idx < len(*items); idx++ {
		if idx > retNum {
			break
		}
		itemInfo := &faiss_index.ItemInfo{
			ItemId: (*items)[idx],
			Score:  (*scores)[idx],
		}
		rankResult = append(rankResult, itemInfo)
	}
	logs.Debug(requestId, time.Now(), "rank result:", inferSample)

	//format result.
	rankRst, err := d.baseModel.InferResultFormat(&rankResult, exposureList)
	if err != nil {
		logs.Error(requestId, time.Now(), err)
		return nil, err
	}
	response["data"] = *rankRst
	logs.Debug(requestId, time.Now(), "format result", response)

	d.baseModel.GetBigCacheRsp().Set(cacheKeyPrefix, []byte(utils.ConvertStructToJson(response)))

	return response, nil
}

func (d *DIN) RecallSampleDirector(userId string, offlineFeature bool, onlineFeature bool, featureList []string) (feature.ExampleFeatures, error) {
	//panic("no not need to Implement interface methods")
	sample := feature.ExampleFeatures{}
	return sample, nil
}

func (d *DIN) RankingSampleDirector(userId string, offlineFeature bool, onlineFeature bool, itemIdList []string, featureList []string) (feature.ExampleFeatures, error) {
	//panic("no not need to Implement interface methods")
	sample := feature.ExampleFeatures{}
	return sample, nil
}
