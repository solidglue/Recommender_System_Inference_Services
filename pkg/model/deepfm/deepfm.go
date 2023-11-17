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
	basemodel.BaseModel // extend baseModel
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

func (d *DeepFM) ModelInferSkywalking(requestId string, userId string, itemList []string, r *http.Request) (map[string]interface{}, error) {
	response := make(map[string]interface{}, 0)
	tensorName := "scores"
	cacheKeyPrefix := userId + d.BaseModel.GetServiceConfig().GetServiceId() + d.BaseModel.GetModelName()

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
	spanUnionEmFv, _, err := internal.Tracer.CreateLocalSpan(r.Context())
	if err != nil {
		return nil, err
	}
	spanUnionEmFv.SetOperationName("get rank infer examples func")
	spanUnionEmFv.Log(time.Now())
	examples, err := d.GetInferExampleFeatures(userId, itemList)
	if err != nil {
		return nil, err
	}
	spanUnionEmFv.Log(time.Now())
	spanUnionEmFv.End()
	logs.Debug(requestId, time.Now(), "examples", examples)

	// get rank scores from tfserving model.
	rankResult := make([]*faiss_index.ItemInfo, 0)
	spanUnionEmFv, _, err = internal.Tracer.CreateLocalSpan(r.Context())
	if err != nil {
		return nil, err
	}
	spanUnionEmFv.SetOperationName("get rank scores func")
	spanUnionEmFv.Log(time.Now())
	items, scores, err := d.rankPredict(examples, tensorName) // d.getServiceConfig().GetModelConfig().rankPredict(examples, tensorName)
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
	spanUnionEmOut, _, err := internal.Tracer.CreateLocalSpan(r.Context())
	if err != nil {
		return nil, err
	}
	spanUnionEmOut.SetOperationName("get rank result func")
	spanUnionEmOut.Log(time.Now())
	rankRst, err := d.BaseModel.InferResultFormat(&rankResult)
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

func (d *DeepFM) ModelInferNoSkywalking(requestId string, userId string, itemList []string, r *http.Request) (map[string]interface{}, error) {
	response := make(map[string]interface{}, 0)
	tensorName := "scores"
	cacheKeyPrefix := userId + d.BaseModel.GetServiceConfig().GetServiceId() + d.BaseModel.GetModelName()

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
	examples, err := d.GetInferExampleFeatures(userId, itemList)
	if err != nil {
		logs.Error(requestId, time.Now(), err)
		return nil, err
	}
	logs.Debug(requestId, time.Now(), "examples", examples)

	// get rank scores from tfserving model.
	rankResult := make([]*faiss_index.ItemInfo, 0)
	items, scores, err := d.rankPredict(examples, tensorName) //d.getServiceConfig().GetModelConfig().rankPredict(examples, tensorName)
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
	rankRst, err := d.BaseModel.InferResultFormat(&rankResult)
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
	scores, err := d.BaseModel.RequestTfservering(&userExamples, &userContextExamples, &itemExamples, tensorName)

	if err != nil {
		return nil, nil, err
	}

	return &items, scores, nil
}

// @overwirte
func (d *DeepFM) GetInferExampleFeatures(userId string, itemList []string) (internal.ExampleFeatures, error) {
	cacheKeyPrefix := userId + d.BaseModel.GetServiceConfig().GetServiceId() + d.BaseModel.GetModelName() + "_samples"

	//init examples
	userExampleFeatures := &internal.SeqExampleBuff{}
	userContextExampleFeatures := &internal.SeqExampleBuff{}
	itemExampleFeaturesList := make([]internal.SeqExampleBuff, 0)
	exampleData := internal.ExampleFeatures{
		UserExampleFeatures:        userExampleFeatures,
		UserContextExampleFeatures: userContextExampleFeatures,
		ItemSeqExampleFeatures:     &itemExampleFeaturesList,
	}

	//set cache
	bigCache, err := bigcache.NewBigCache(bigCacheConfDeepfm)
	if err != nil {
		return exampleData, err
	}

	// if hit cache.
	if lifeWindowS > 0 {
		exampleDataBytes, _ := bigCache.Get(cacheKeyPrefix)
		err = json.Unmarshal(exampleDataBytes, &exampleData)
		if err != nil {
			return exampleData, err
		}
		return exampleData, nil

	}

	// if not hit cache, get features from redis and request.
	userExampleFeatures, err = d.BaseModel.GetUserExampleFeatures(userId)
	if err != nil {
		return exampleData, err
	}
	userContextExampleFeatures, err = d.BaseModel.GetUserContextExampleFeatures(userId)
	if err != nil {
		return exampleData, err
	}

	//get items features.
	itemExampleFeaturesTmp, err := d.getItemExamplesFeatures(itemList)
	if err != nil {
		return exampleData, err
	}

	itemExampleFeaturesList = *itemExampleFeaturesTmp
	exampleData = internal.ExampleFeatures{
		UserExampleFeatures:        userExampleFeatures,
		UserContextExampleFeatures: userContextExampleFeatures,
		ItemSeqExampleFeatures:     &itemExampleFeaturesList,
	}

	if lifeWindowS > 0 {
		bigCache.Set(cacheKeyPrefix, []byte(utils.ConvertStructToJson(exampleData)))
	}

	return exampleData, nil
}

func (d *DeepFM) getItemExamplesFeatures(itemList []string) (*[]internal.SeqExampleBuff, error) {
	//TODO: use bloom filter check items, avoid all items search redis.
	redisKeyPrefix := d.BaseModel.GetServiceConfig().GetModelConfig().GetItemRedisKeyPre()
	itemSeqExampleBuffs := make([]internal.SeqExampleBuff, 0)
	for _, itemId := range itemList {
		redisKey := redisKeyPrefix + itemId
		if d.BaseModel.GetItemBloomFilter().Test([]byte(itemId)) {
			userExampleFeats, err := d.BaseModel.GetServiceConfig().GetRedisConfig().GetRedisPool().Get(redisKey)
			itemExampleFeatsBuff := make([]byte, 0)
			if err != nil {
				return &itemSeqExampleBuffs, nil
			} else {
				itemExampleFeatsBuff = []byte(userExampleFeats)
			}

			itemSeqExampleBuff := internal.SeqExampleBuff{
				Key:  &itemId,
				Buff: &itemExampleFeatsBuff,
			}
			itemSeqExampleBuffs = append(itemSeqExampleBuffs, itemSeqExampleBuff)
		}
	}

	return &itemSeqExampleBuffs, nil
}
