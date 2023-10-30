package dssm

import (
	"encoding/json"
	"infer-microservices/common"
	"infer-microservices/common/flags"
	"infer-microservices/utils"
	"infer-microservices/utils/logs"
	"time"

	"github.com/allegro/bigcache"
)

var bigCacheConfRecallSamples bigcache.Config
var bigCacheConfRankSamples bigcache.Config
var lifeWindowS time.Duration
var cleanWindowS time.Duration
var hardMaxCacheSize int
var maxEntrySize int
var maxEntriesInWindow int
var verbose bool
var shards int

func init() {
	flagFactory := flags.FlagFactory{}
	flagCache := flagFactory.CreateFlagCache()

	lifeWindowS = time.Duration(*flagCache.GetBigcacheLifeWindowS())
	cleanWindowS = time.Duration(*flagCache.GetBigcacheCleanWindowS())
	hardMaxCacheSize = *flagCache.GetBigcacheHardMaxCacheSize()
	maxEntrySize = *flagCache.GetBigcacheMaxEntrySize()

	bigCacheConfRecallSamples = bigcache.Config{
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

func (d *Dssm) GetInferExampleFeatures() (common.ExampleFeatures, error) {
	cacheKeyPrefix := d.getUserId() + d.serviceConfig.GetServiceId() + "_recallSamples"

	//init examples
	userExampleFeatures := &common.SeqExampleBuff{}
	userContextExampleFeatures := &common.SeqExampleBuff{}
	exampleData := common.ExampleFeatures{
		UserExampleFeatures:        userExampleFeatures,
		UserContextExampleFeatures: userContextExampleFeatures,
	}

	//set cache
	bigCache, err := bigcache.NewBigCache(bigCacheConfRecallSamples)
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
	userExampleFeatures, err = d.getUserExampleFeatures()
	if err != nil {
		logs.Error(err)
		return exampleData, err
	}
	userContextExampleFeatures, err = d.getUserContextExampleFeatures()
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
		bigCache.Set(cacheKeyPrefix, []byte(utils.Struct2Json(exampleData)))
	}

	return exampleData, nil
}

func (d *Dssm) getUserExampleFeatures() (*common.SeqExampleBuff, error) {

	//TODO: update context features.
	redisKey := d.getServiceConfig().GetModelConfig().GetUserRedisKeyPre() + d.getUserId()
	userExampleFeats, err := d.getServiceConfig().GetRedisConfig().GetRedisPool().Get(redisKey)

	userSeqExampleBuff := common.SeqExampleBuff{}
	userExampleFeatsBuff := make([]byte, 0)

	if err != nil {
		logs.Error("get item features err", err)
		return &userSeqExampleBuff, err
	} else {
		userExampleFeatsBuff = []byte(userExampleFeats) //.(string)
	}

	//protrait features & realtime features.
	userSeqExampleBuff = common.SeqExampleBuff{
		Key:  &d.userId,
		Buff: &userExampleFeatsBuff,
	}

	return &userSeqExampleBuff, nil
}

func (d *Dssm) getUserContextExampleFeatures() (*common.SeqExampleBuff, error) {
	userContextSeqExampleBuff := common.SeqExampleBuff{}
	userContextExampleFeatsBuff := make([]byte, 0)

	//TODO: update context features. only from requst. such as location , time
	//context features.
	userContextSeqExampleBuff = common.SeqExampleBuff{
		Key:  &d.userId,
		Buff: &userContextExampleFeatsBuff,
	}

	return &userContextSeqExampleBuff, nil
}
