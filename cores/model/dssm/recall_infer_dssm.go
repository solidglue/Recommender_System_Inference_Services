package dssm

import (
	"encoding/json"
	"infer-microservices/common"
	faiss_index "infer-microservices/common/faiss_gogofaster"
	"infer-microservices/cores/service_config"
	"infer-microservices/utils"

	"infer-microservices/utils/logs"
	"net/http"
	"sync"
	"time"

	"github.com/allegro/bigcache"
)

var bigCacheConfRecallResult bigcache.Config

type Dssm struct {
	userId        string
	retNum        int
	serviceConfig *service_config.ServiceConfig
}

func init() {
	bigCacheConfRecallResult = bigcache.Config{
		Shards:             shards,
		LifeWindow:         lifeWindowS * time.Minute,
		CleanWindow:        cleanWindowS * time.Minute,
		MaxEntriesInWindow: maxEntriesInWindow,
		MaxEntrySize:       maxEntrySize,
		Verbose:            verbose,
		HardMaxCacheSize:   hardMaxCacheSize,
		OnRemove:           nil,
		OnRemoveWithReason: nil,
	}
}

// userid
func (d *Dssm) SetUserId(userId string) {
	d.userId = userId
}

func (d *Dssm) getUserId() string {
	return d.userId
}

// retNum
func (d *Dssm) SetRetNum(retNum int) {
	d.retNum = retNum
}

func (d *Dssm) getRetNum() int {
	return d.retNum
}

// serviceConfig *service_config.ServiceConfig
func (d *Dssm) SetServiceConfig(serviceConfig *service_config.ServiceConfig) {
	d.serviceConfig = serviceConfig
}

func (d *Dssm) getServiceConfig() *service_config.ServiceConfig {
	return d.serviceConfig
}

func (d *Dssm) RecallInferSkywalking(r *http.Request) (map[string]interface{}, error) {
	response := make(map[string]interface{}, 0)
	cacheKeyPrefix := d.getUserId() + d.serviceConfig.GetServiceId() + "_recallResult"
	tensorName := "user_embedding"

	//set cache
	bigCache, err := bigcache.NewBigCache(bigCacheConfRecallResult)
	if err != nil {
		logs.Error(err)
	}

	// get features from cache.
	if lifeWindowS > 0 {
		exampleDataBytes, _ := bigCache.Get(cacheKeyPrefix)
		err = json.Unmarshal(exampleDataBytes, &response)
		if err != nil {
			logs.Error(err)
		}
		return response, nil
	}

	//get infer samples.
	examples := common.ExampleFeatures{}
	spanUnionEmFv, _, err := common.Tracer.CreateLocalSpan(r.Context())
	if err != nil {
		return nil, err
	}

	spanUnionEmFv.SetOperationName("get recall infer examples func")
	spanUnionEmFv.Log(time.Now())
	examples, err = d.getInferExampleFeatures()
	if err != nil {
		return nil, err
	}
	spanUnionEmFv.Log(time.Now())
	spanUnionEmFv.End()

	// get embedding from tfserving model.
	embeddingVector := make([]float32, 0)
	spanUnionEmFv, _, err = common.Tracer.CreateLocalSpan(r.Context())
	if err != nil {
		return nil, err
	}
	spanUnionEmFv.SetOperationName("get recall embedding func")
	spanUnionEmFv.Log(time.Now())

	embeddingVector_, err := d.getServiceConfig().GetModelClient().Embedding(examples, tensorName)
	if err != nil {
		logs.Error(err)
		return nil, err
	} else {
		embeddingVector = *embeddingVector_
	}
	spanUnionEmFv.Log(time.Now())
	spanUnionEmFv.End()

	//request faiss index
	recallResult := make([]*faiss_index.ItemInfo, 0)
	spanUnionEmFr, _, err1 := common.Tracer.CreateLocalSpan(r.Context())
	if err1 != nil {
		logs.Error(err)
		return nil, err1
	}
	spanUnionEmFr.SetOperationName("get recall faiss index func")
	spanUnionEmFr.Log(time.Now())
	recallResult, err = d.getServiceConfig().GetFaissIndexClient().FaissVectorSearch(examples, embeddingVector)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	spanUnionEmFr.Log(time.Now())
	spanUnionEmFr.End()

	//format result.
	spanUnionEmOut, _, err := common.Tracer.CreateLocalSpan(r.Context())
	if err != nil {
		return nil, err
	}
	spanUnionEmOut.SetOperationName("get recall result func")

	spanUnionEmOut.Log(time.Now())
	recallRst, err := d.recallResultFmt(&recallResult)
	if err != nil {
		return nil, err
	}
	spanUnionEmOut.Log(time.Now())
	spanUnionEmOut.End()

	if len(*recallRst) == 0 {
		logs.Error("recall 0 item, check the faiss index plz. ")
		return nil, err
	}

	response["data"] = *recallRst

	if lifeWindowS > 0 {
		bigCache.Set(cacheKeyPrefix, []byte(utils.Struct2Json(response)))
	}

	return response, nil
}

func (d *Dssm) RecallInferNoSkywalking(r *http.Request) (map[string]interface{}, error) {
	response := make(map[string]interface{}, 0)
	cacheKeyPrefix := d.getUserId() + d.serviceConfig.GetServiceId() + "_recallResult"
	tensorName := "user_embedding"

	//set cache
	bigCache, err := bigcache.NewBigCache(bigCacheConfRecallResult)
	if err != nil {
		return nil, err
	}

	// get features from cache.
	if lifeWindowS > 0 {
		exampleDataBytes, _ := bigCache.Get(cacheKeyPrefix)
		err = json.Unmarshal(exampleDataBytes, &response)
		if err != nil {
			return nil, err
		}
		return response, nil
	}

	//get infer samples.
	examples := common.ExampleFeatures{}
	examples, err = d.getInferExampleFeatures()

	// get embedding from tfserving model.
	embeddingVector := make([]float32, 0)
	embeddingVector_, err := d.getServiceConfig().GetModelClient().Embedding(examples, tensorName)
	if err != nil {
		return nil, err
	} else {
		embeddingVector = *embeddingVector_
	}

	//request faiss index
	recallResult := make([]*faiss_index.ItemInfo, 0)
	recallResult, err = d.getServiceConfig().GetFaissIndexClient().FaissVectorSearch(examples, embeddingVector)
	if err != nil {
		return nil, err
	}

	//format result.
	recallRst, err := d.recallResultFmt(&recallResult)
	if err != nil {
		return nil, err
	}

	if len(*recallRst) == 0 {
		logs.Error("recall 0 item, check the faiss index plz. ")
		return nil, err
	}

	response["data"] = *recallRst

	if lifeWindowS > 0 {
		bigCache.Set(cacheKeyPrefix, []byte(utils.Struct2Json(response)))
	}

	return response, nil
}

func (d *Dssm) recallResultFmt(recallResult *[]*faiss_index.ItemInfo) (*[]map[string]interface{}, error) {
	recall := make([]map[string]interface{}, 0)
	recallTmp := make(chan map[string]interface{}, len(*recallResult)) // 20221011
	var wg sync.WaitGroup

	for idx := 0; idx < len(*recallResult); idx++ {
		rawCell := (*recallResult)[idx]
		wg.Add(1)
		go func(raw_cell_ *faiss_index.ItemInfo) {
			defer wg.Done()
			returnCell := make(map[string]interface{})
			returnCell["itemid"] = raw_cell_.ItemId
			returnCell["score"] = utils.FloatRound(raw_cell_.Score, 4)
			recallTmp <- returnCell
		}(rawCell)
	}
	wg.Wait()
	for idx := 0; idx < len(*recallResult); idx++ {
		returnCellTmp := <-recallTmp
		recall = append(recall, returnCellTmp)
	}
	close(recallTmp)

	return &recall, nil
}