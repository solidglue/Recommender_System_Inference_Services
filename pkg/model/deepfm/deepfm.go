package deepfm

import (
	"encoding/json"
	"infer-microservices/internal"
	faiss_index "infer-microservices/internal/faiss_gogofaster"
	"infer-microservices/internal/flags"
	"infer-microservices/pkg/logs"
	"infer-microservices/pkg/model/basemodel"
	"infer-microservices/pkg/utils"
	"net/http"
	"time"

	"github.com/allegro/bigcache"
)

var bigCacheConfDeepfm bigcache.Config
var lifeWindowS time.Duration
var cleanWindowS time.Duration
var hardMaxCacheSize int
var maxEntrySize int
var maxEntriesInWindow int
var verbose bool
var shards int

type DeepFM struct {
	basemodel basemodel.BaseModel // extend baseModel
	modelType string
}

func init() {
	flagFactory := flags.FlagFactory{}
	flagCache := flagFactory.CreateFlagCache()
	lifeWindowS = time.Duration(*flagCache.GetBigcacheLifeWindowS())
	cleanWindowS = time.Duration(*flagCache.GetBigcacheCleanWindowS())
	hardMaxCacheSize = *flagCache.GetBigcacheHardMaxCacheSize()
	maxEntrySize = *flagCache.GetBigcacheMaxEntrySize()
	bigCacheConfDeepfm = bigcache.Config{
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

func (d *DeepFM) SetBaseModel(basemodel basemodel.BaseModel) {
	d.basemodel = basemodel
}

// modeltype
func (d *DeepFM) SetModelType(modelType string) {
	d.modelType = modelType
}

func (d *DeepFM) GetModelType() string {
	return d.modelType
}

func (d *DeepFM) ModelInferSkywalking(requestId string, userId string, itemList []string, r *http.Request, createSample basemodel.CreateSampleCallBackFunc) (map[string]interface{}, error) {
	response := make(map[string]interface{}, 0)
	tensorName := "scores"
	cacheKeyPrefix := userId + d.basemodel.GetServiceConfig().GetServiceId() + d.basemodel.GetModelName()

	//set cache
	bigCache, err := bigcache.NewBigCache(bigCacheConfDeepfm)
	if err != nil {
		logs.Error(requestId, time.Now(), err)
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
	spanUnionEmFv, _, err := internal.GetTracer().CreateLocalSpan(r.Context())
	if err != nil {
		return nil, err
	}
	spanUnionEmFv.SetOperationName("get rank infer examples func")
	spanUnionEmFv.Log(time.Now())
	examples, err := createSample(userId, itemList) //create sample by callback func
	if err != nil {
		return nil, err
	}
	spanUnionEmFv.Log(time.Now())
	spanUnionEmFv.End()
	logs.Debug(requestId, time.Now(), "examples", examples)

	// get rank scores from tfserving model.
	rankResult := make([]*faiss_index.ItemInfo, 0)
	spanUnionEmFv, _, err = internal.GetTracer().CreateLocalSpan(r.Context())
	if err != nil {
		return nil, err
	}
	spanUnionEmFv.SetOperationName("get rank scores func")
	spanUnionEmFv.Log(time.Now())
	items, scores, err := d.rankPredict(examples, tensorName)
	if err != nil {
		return nil, err
	}
	spanUnionEmFv.Log(time.Now())
	spanUnionEmFv.End()
	logs.Debug(requestId, time.Now(), "unrank result:", examples)

	//build rank result whith tfserving.ItemInfo
	for idx := 0; idx < len(*items); idx++ {
		itemInfo := &faiss_index.ItemInfo{
			ItemId: (*items)[idx],
			Score:  (*scores)[idx],
		}
		rankResult = append(rankResult, itemInfo)
	}
	logs.Debug(requestId, time.Now(), "rank result:", examples)

	//format result.
	spanUnionEmOut, _, err := internal.GetTracer().CreateLocalSpan(r.Context())
	if err != nil {
		return nil, err
	}
	spanUnionEmOut.SetOperationName("get rank result func")
	spanUnionEmOut.Log(time.Now())
	rankRst, err := d.basemodel.InferResultFormat(&rankResult)
	if err != nil {
		return nil, err
	}
	spanUnionEmOut.Log(time.Now())
	spanUnionEmOut.End()

	response["data"] = *rankRst
	logs.Debug(requestId, time.Now(), "format result", response)

	if lifeWindowS > 0 {
		bigCache.Set(cacheKeyPrefix, []byte(utils.ConvertStructToJson(response)))
	}

	return response, nil
}

func (d *DeepFM) ModelInferNoSkywalking(requestId string, userId string, itemList []string, r *http.Request, createSample basemodel.CreateSampleCallBackFunc) (map[string]interface{}, error) {
	response := make(map[string]interface{}, 0)
	tensorName := "scores"
	cacheKeyPrefix := userId + d.basemodel.GetServiceConfig().GetServiceId() + d.basemodel.GetModelName()

	//set cache
	bigCache, err := bigcache.NewBigCache(bigCacheConfDeepfm)
	if err != nil {
		logs.Error(requestId, time.Now(), err)
	}

	// get features from cache.
	if lifeWindowS > 0 {
		exampleDataBytes, _ := bigCache.Get(cacheKeyPrefix)
		err = json.Unmarshal(exampleDataBytes, &response)
		if err != nil {
			logs.Error(requestId, time.Now(), err)
			return nil, err
		}
		return response, nil

	}

	//get infer samples.
	examples, err := createSample(userId, itemList) //create sample by callback func
	if err != nil {
		logs.Error(requestId, time.Now(), err)
		return nil, err
	}
	logs.Debug(requestId, time.Now(), "examples", examples)

	// get rank scores from tfserving model.
	rankResult := make([]*faiss_index.ItemInfo, 0)
	items, scores, err := d.rankPredict(examples, tensorName)
	if err != nil {
		logs.Error(requestId, time.Now(), err)
		return nil, err
	}
	logs.Debug(requestId, time.Now(), "unrank result:", examples)

	//build rank result whith tfserving.ItemInfo
	for idx := 0; idx < len(*items); idx++ {
		itemInfo := &faiss_index.ItemInfo{
			ItemId: (*items)[idx],
			Score:  (*scores)[idx],
		}
		rankResult = append(rankResult, itemInfo)
	}
	logs.Debug(requestId, time.Now(), "rank result:", examples)

	//format result.
	rankRst, err := d.basemodel.InferResultFormat(&rankResult)
	if err != nil {
		logs.Error(requestId, time.Now(), err)
		return nil, err
	}
	response["data"] = *rankRst
	logs.Debug(requestId, time.Now(), "format result", response)

	if lifeWindowS > 0 {
		bigCache.Set(cacheKeyPrefix, []byte(utils.ConvertStructToJson(response)))
	}

	return response, nil
}

// request rank scores from tfserving
func (d *DeepFM) rankPredict(examples internal.ExampleFeatures, tensorName string) (*[]string, *[]float32, error) {

	userExamples := make([][]byte, 0)
	userContextExamples := make([][]byte, 0)
	itemExamples := make([][]byte, 0)
	items := make([]string, 0)

	userExamples = append(userExamples, *(examples.UserExampleFeatures.Buff))
	userContextExamples = append(userContextExamples, *(examples.UserContextExampleFeatures.Buff))

	for _, itemExample := range *examples.ItemSeqExampleFeatures {
		items = append(items, *(itemExample.Key))
		itemExamples = append(itemExamples, *(itemExample.Buff))
	}
	scores, err := d.basemodel.RequestTfservering(&userExamples, &userContextExamples, &itemExamples, tensorName)

	if err != nil {
		return nil, nil, err
	}

	return &items, scores, nil
}
