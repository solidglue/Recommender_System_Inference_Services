package dssm

import (
	"encoding/json"
	"infer-microservices/common"
	faiss_index "infer-microservices/common/faiss_gogofaster"
	"infer-microservices/common/flags"
	"infer-microservices/cores/faiss"
	"infer-microservices/cores/model/basemodel"
	"infer-microservices/utils"
	"infer-microservices/utils/logs"
	"net/http"
	"time"

	"github.com/allegro/bigcache"
)

var bigCacheConfDssm bigcache.Config
var lifeWindowS time.Duration
var cleanWindowS time.Duration
var hardMaxCacheSize int
var maxEntrySize int
var maxEntriesInWindow int
var verbose bool
var shards int

type Dssm struct {
	basemodel.BaseModel // extend baseModel
	retNum              int
}

func init() {
	flagFactory := flags.FlagFactory{}
	flagCache := flagFactory.CreateFlagCache()
	lifeWindowS = time.Duration(*flagCache.GetBigcacheLifeWindowS())
	cleanWindowS = time.Duration(*flagCache.GetBigcacheCleanWindowS())
	hardMaxCacheSize = *flagCache.GetBigcacheHardMaxCacheSize()
	maxEntrySize = *flagCache.GetBigcacheMaxEntrySize()
	bigCacheConfDssm = bigcache.Config{
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

// retNum
func (d *Dssm) SetRetNum(retNum int) {
	d.retNum = retNum
}

func (d *Dssm) GetRetNum() int {
	return d.retNum
}

func (d *Dssm) ModelInferSkywalking(r *http.Request) (map[string]interface{}, error) {
	response := make(map[string]interface{}, 0)
	cacheKeyPrefix := d.BaseModel.GetUserId() + d.BaseModel.GetServiceConfig().GetServiceId() + d.BaseModel.GetModelName()

	tensorName := "user_embedding"

	//set cache
	bigCache, err := bigcache.NewBigCache(bigCacheConfDssm)
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
	spanUnionEmFv, _, err := common.Tracer.CreateLocalSpan(r.Context())
	if err != nil {
		return nil, err
	}

	spanUnionEmFv.SetOperationName("get recall infer examples func")
	spanUnionEmFv.Log(time.Now())
	examples, err := d.GetInferExampleFeatures()
	if err != nil {
		return nil, err
	}
	spanUnionEmFv.Log(time.Now())
	spanUnionEmFv.End()

	// get embedding from tfserving model.
	spanUnionEmFv, _, err = common.Tracer.CreateLocalSpan(r.Context())
	if err != nil {
		return nil, err
	}
	spanUnionEmFv.SetOperationName("get recall embedding func")
	spanUnionEmFv.Log(time.Now())

	embeddingVector, err := d.embedding(examples, tensorName)
	if err != nil {
		logs.Error(err)
		return nil, err
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
	recallResult, err = faiss.FaissVectorSearch(d.BaseModel.GetServiceConfig().GetFaissIndexConfig(), examples, *embeddingVector)
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
	recallRst, err := d.BaseModel.InferResultFormat(&recallResult)
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
		bigCache.Set(cacheKeyPrefix, []byte(utils.ConvertStructToJson(response)))
	}

	return response, nil
}

func (d *Dssm) ModelInferNoSkywalking(r *http.Request) (map[string]interface{}, error) {
	response := make(map[string]interface{}, 0)
	cacheKeyPrefix := d.BaseModel.GetUserId() + d.BaseModel.GetServiceConfig().GetServiceId() + d.BaseModel.GetModelName()
	tensorName := "user_embedding"

	//set cache
	bigCache, err := bigcache.NewBigCache(bigCacheConfDssm)
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
	examples, err := d.GetInferExampleFeatures()
	if err != nil {
		return nil, err
	}

	// get embedding from tfserving model.
	embeddingVector, err := d.embedding(examples, tensorName)
	if err != nil {
		return nil, err
	}

	//request faiss index
	recallResult := make([]*faiss_index.ItemInfo, 0)
	recallResult, err = faiss.FaissVectorSearch(d.BaseModel.GetServiceConfig().GetFaissIndexConfig(), examples, *embeddingVector)
	if err != nil {
		return nil, err
	}

	//format result.
	recallRst, err := d.BaseModel.InferResultFormat(&recallResult)
	if err != nil {
		return nil, err
	}

	if len(*recallRst) == 0 {
		logs.Error("recall 0 item, check the faiss index plz. ")
		return nil, err
	}

	response["data"] = *recallRst

	if lifeWindowS > 0 {
		bigCache.Set(cacheKeyPrefix, []byte(utils.ConvertStructToJson(response)))
	}

	return response, nil
}

// request embedding vector from tfserving
func (d *Dssm) embedding(examples common.ExampleFeatures, tensorName string) (*[]float32, error) {

	userExamples := make([][]byte, 0)
	userContextExamples := make([][]byte, 0)
	itemExamples := make([][]byte, 0)

	userExamples = append(userExamples, *(examples.UserExampleFeatures.Buff))
	userContextExamples = append(userContextExamples, *(examples.UserContextExampleFeatures.Buff))

	response, err := d.BaseModel.RequestTfservering(&userExamples, &itemExamples, &userContextExamples, tensorName)
	if err != nil {
		logs.Error(err)
		return nil, err
	}

	return response, nil
}

// @overwrite
func (d *Dssm) GetInferExampleFeatures() (common.ExampleFeatures, error) {
	cacheKeyPrefix := d.BaseModel.GetUserId() + d.BaseModel.GetServiceConfig().GetServiceId() + d.BaseModel.GetModelName() + "_samples"

	//init examples
	userExampleFeatures := &common.SeqExampleBuff{}
	userContextExampleFeatures := &common.SeqExampleBuff{}
	exampleData := common.ExampleFeatures{
		UserExampleFeatures:        userExampleFeatures,
		UserContextExampleFeatures: userContextExampleFeatures,
	}

	//set cache
	bigCache, err := bigcache.NewBigCache(bigCacheConfDssm)
	if err != nil {
		logs.Error(err)
	}

	// if hit cacha.
	if lifeWindowS > 0 {

		//INFO:MMO, go-cache can't set MaxCacheSize. change to use bigcache.

		// if cacheData, ok := goCache.Get(cacheKeyPrefix); ok {
		// 	return cacheData.(ExampleFeatures), nil
		// }

		exampleDataBytes, _ := bigCache.Get(cacheKeyPrefix)
		err = json.Unmarshal(exampleDataBytes, &exampleData)
		if err != nil {
			logs.Error(err)
		}
		return exampleData, nil

	}

	// if not hit cache, get features from redis and request.
	userExampleFeatures, err = d.BaseModel.GetUserExampleFeatures()
	if err != nil {
		logs.Error(err)
		return exampleData, err
	}
	userContextExampleFeatures, err = d.BaseModel.GetUserContextExampleFeatures()
	if err != nil {
		logs.Error(err)
		return exampleData, err
	}

	exampleData = common.ExampleFeatures{
		UserExampleFeatures:        userExampleFeatures,
		UserContextExampleFeatures: userContextExampleFeatures,
	}

	if lifeWindowS > 0 {
		// goCache.Set(cacheKeyPrefix, &exampleData, cacheTimeSecond)
		bigCache.Set(cacheKeyPrefix, []byte(utils.ConvertStructToJson(exampleData)))
	}

	return exampleData, nil
}
