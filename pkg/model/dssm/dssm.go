package dssm

import (
	"encoding/json"
	"infer-microservices/internal"
	faiss_index "infer-microservices/internal/faiss_gogofaster"
	"infer-microservices/internal/flags"
	"infer-microservices/pkg/faiss"
	"infer-microservices/pkg/logs"
	"infer-microservices/pkg/model/basemodel"
	"infer-microservices/pkg/utils"
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
	basemodel basemodel.BaseModel // extend baseModel
	retNum    int
	modelType string
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

func (d *Dssm) SetBaseModel(basemodel basemodel.BaseModel) {
	d.basemodel = basemodel
}

// modeltype
func (d *Dssm) SetModelType(modelType string) {
	d.modelType = modelType
}

func (d *Dssm) GetModelType() string {
	return d.modelType
}

func (d *Dssm) ModelInferSkywalking(requestId string, userId string, itemList []string, r *http.Request, createSample basemodel.CreateSampleCallBackFunc) (map[string]interface{}, error) {
	response := make(map[string]interface{}, 0)
	cacheKeyPrefix := userId + d.basemodel.GetServiceConfig().GetServiceId() + d.basemodel.GetModelName()

	tensorName := "user_embedding"

	//set cache
	bigCache, err := bigcache.NewBigCache(bigCacheConfDssm)
	if err != nil {
		logs.Error(requestId, time.Now(), err)
	}

	// get features from cache.
	if lifeWindowS > 0 {
		exampleDataBytes, _ := bigCache.Get(cacheKeyPrefix)
		err = json.Unmarshal(exampleDataBytes, &response)
		if err != nil {
			logs.Error(requestId, time.Now(), err)
		}
		return response, nil
	}

	//get infer samples.
	spanUnionEmFv, _, err := internal.Tracer.CreateLocalSpan(r.Context())
	if err != nil {
		return nil, err
	}

	spanUnionEmFv.SetOperationName("get recall infer examples func")
	spanUnionEmFv.Log(time.Now())
	examples, err := createSample(userId, itemList) //create sample by callback func
	if err != nil {
		return nil, err
	}
	spanUnionEmFv.Log(time.Now())
	spanUnionEmFv.End()
	logs.Debug(requestId, time.Now(), "example:", examples)

	// get embedding from tfserving model.
	spanUnionEmFv, _, err = internal.Tracer.CreateLocalSpan(r.Context())
	if err != nil {
		return nil, err
	}
	spanUnionEmFv.SetOperationName("get recall embedding func")
	spanUnionEmFv.Log(time.Now())

	embeddingVector, err := d.embedding(examples, tensorName)
	if err != nil {
		logs.Error(requestId, time.Now(), err)
		return nil, err
	}
	spanUnionEmFv.Log(time.Now())
	spanUnionEmFv.End()
	logs.Debug(requestId, time.Now(), "embeddingVector:", embeddingVector)

	//request faiss index
	recallResult := make([]*faiss_index.ItemInfo, 0)
	spanUnionEmFr, _, err1 := internal.Tracer.CreateLocalSpan(r.Context())
	if err1 != nil {
		logs.Error(requestId, time.Now(), err)
		return nil, err1
	}
	spanUnionEmFr.SetOperationName("get recall faiss index func")
	spanUnionEmFr.Log(time.Now())
	recallResult, err = faiss.FaissVectorSearch(d.basemodel.GetServiceConfig().GetFaissIndexConfig(), examples, *embeddingVector)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	spanUnionEmFr.Log(time.Now())
	spanUnionEmFr.End()
	logs.Debug(requestId, time.Now(), "recall result:", recallResult)

	//format result.
	spanUnionEmOut, _, err := internal.Tracer.CreateLocalSpan(r.Context())
	if err != nil {
		return nil, err
	}
	spanUnionEmOut.SetOperationName("get recall result func")

	spanUnionEmOut.Log(time.Now())
	recallRst, err := d.basemodel.InferResultFormat(&recallResult)
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
	logs.Debug(requestId, time.Now(), "format result:", recallResult)

	if lifeWindowS > 0 {
		bigCache.Set(cacheKeyPrefix, []byte(utils.ConvertStructToJson(response)))
	}

	return response, nil
}

func (d *Dssm) ModelInferNoSkywalking(requestId string, userId string, itemList []string, r *http.Request, createSample basemodel.CreateSampleCallBackFunc) (map[string]interface{}, error) {
	response := make(map[string]interface{}, 0)
	cacheKeyPrefix := userId + d.basemodel.GetServiceConfig().GetServiceId() + d.basemodel.GetModelName()
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
	examples, err := createSample(userId, itemList) //create sample by callback func
	if err != nil {
		return nil, err
	}

	// get embedding from tfserving model.
	embeddingVector, err := d.embedding(examples, tensorName)
	if err != nil {
		return nil, err
	}
	logs.Debug(requestId, time.Now(), "embeddingVector:", embeddingVector)

	//request faiss index
	recallResult := make([]*faiss_index.ItemInfo, 0)
	recallResult, err = faiss.FaissVectorSearch(d.basemodel.GetServiceConfig().GetFaissIndexConfig(), examples, *embeddingVector)
	if err != nil {
		return nil, err
	}
	logs.Debug(requestId, time.Now(), "recall result:", recallResult)

	//format result.
	recallRst, err := d.basemodel.InferResultFormat(&recallResult)
	if err != nil {
		return nil, err
	}

	if len(*recallRst) == 0 {
		logs.Error("recall 0 item, check the faiss index plz. ")
		return nil, err
	}
	response["data"] = *recallRst
	logs.Debug(requestId, time.Now(), "format result:", recallResult)

	if lifeWindowS > 0 {
		bigCache.Set(cacheKeyPrefix, []byte(utils.ConvertStructToJson(response)))
	}

	return response, nil
}

// request embedding vector from tfserving
func (d *Dssm) embedding(examples internal.ExampleFeatures, tensorName string) (*[]float32, error) {

	userExamples := make([][]byte, 0)
	userContextExamples := make([][]byte, 0)
	itemExamples := make([][]byte, 0)

	userExamples = append(userExamples, *(examples.UserExampleFeatures.Buff))
	userContextExamples = append(userContextExamples, *(examples.UserContextExampleFeatures.Buff))

	response, err := d.basemodel.RequestTfservering(&userExamples, &itemExamples, &userContextExamples, tensorName)
	if err != nil {
		logs.Error(err)
		return nil, err
	}

	return response, nil
}
