package deepfm

import (
	"encoding/json"
	"infer-microservices/common"
	"infer-microservices/common/flags"
	"infer-microservices/utils"
	"time"

	"github.com/allegro/bigcache"
)

var lifeWindowS1 time.Duration
var cleanWindowS1 time.Duration
var hardMaxCacheSize1 int
var maxEntrySize1 int
var maxEntriesInWindow1 int
var verbose1 bool
var shards1 int
var bigCacheConfRankSamples bigcache.Config

func init() {
	flagFactory := flags.FlagFactory{}
	flagCache := flagFactory.FlagCacheFactory()

	lifeWindowS1 = time.Duration(*flagCache.GetBigcacheLifeWindowS())
	cleanWindowS1 = time.Duration(*flagCache.GetBigcacheCleanWindowS())
	hardMaxCacheSize1 = *flagCache.GetBigcacheHardMaxCacheSize()
	maxEntrySize1 = *flagCache.GetBigcacheMaxEntrySize()

	bigCacheConfRankSamples = bigcache.Config{
		Shards:             shards1,
		LifeWindow:         lifeWindowS1 * time.Minute,
		CleanWindow:        cleanWindowS1 * time.Minute,
		MaxEntriesInWindow: maxEntriesInWindow1,
		MaxEntrySize:       maxEntrySize1,
		Verbose:            verbose1,
		HardMaxCacheSize:   hardMaxCacheSize1,
		OnRemove:           nil,
		OnRemoveWithReason: nil,
	}
}

func (d *DeepFM) GetInferExampleFeatures() (common.ExampleFeatures, error) {
	cacheKeyPrefix := d.getUserId() + d.serviceConfig.GetServiceId() + "_rankSamples"

	//init examples
	userExampleFeatures := &common.SeqExampleBuff{}
	userContextExampleFeatures := &common.SeqExampleBuff{}
	itemExampleFeaturesList := make([]common.SeqExampleBuff, 0)
	exampleData := common.ExampleFeatures{
		UserExampleFeatures:        userExampleFeatures,
		UserContextExampleFeatures: userContextExampleFeatures,
		ItemSeqExampleFeatures:     &itemExampleFeaturesList,
	}

	//set cache
	bigCache, err := bigcache.NewBigCache(bigCacheConfRankSamples)
	if err != nil {
		return exampleData, err
	}

	// if hit cache.
	if lifeWindowS1 > 0 {
		exampleDataBytes, _ := bigCache.Get(cacheKeyPrefix)
		err = json.Unmarshal(exampleDataBytes, &exampleData)
		if err != nil {
			return exampleData, err
		}
		return exampleData, nil

	}

	// if not hit cache, get features from redis and request.
	userExampleFeatures, err = d.getUserExampleFeatures()
	if err != nil {
		return exampleData, err
	}
	userContextExampleFeatures, err = d.getUserContextExampleFeatures()
	if err != nil {
		return exampleData, err
	}

	//get items features.
	itemExampleFeaturesTmp, err := d.getItemExamplesFeatures()
	if err != nil {
		return exampleData, err
	}

	itemExampleFeaturesList = *itemExampleFeaturesTmp
	exampleData = common.ExampleFeatures{
		UserExampleFeatures:        userExampleFeatures,
		UserContextExampleFeatures: userContextExampleFeatures,
		ItemSeqExampleFeatures:     &itemExampleFeaturesList,
	}

	if lifeWindowS1 > 0 {
		bigCache.Set(cacheKeyPrefix, []byte(utils.Struct2Json(exampleData)))
	}

	return exampleData, nil
}

func (d *DeepFM) getUserExampleFeatures() (*common.SeqExampleBuff, error) {

	redisKey := d.getServiceConfig().GetModelClient().GetUserRedisKeyPre() + d.getUserId()
	userExampleFeatsBuff := make([]byte, 0)
	userSeqExampleBuff := common.SeqExampleBuff{}
	userExampleFeats, err := d.getServiceConfig().GetRedisClient().GetRedisPool().Get(redisKey)
	if err != nil {
		return &userSeqExampleBuff, err
	} else {
		userExampleFeatsBuff = []byte(userExampleFeats)
	}

	//protrait features & realtime features.
	userSeqExampleBuff = common.SeqExampleBuff{
		Key:  &d.userId,
		Buff: &userExampleFeatsBuff,
	}

	return &userSeqExampleBuff, nil
}

func (d *DeepFM) getUserContextExampleFeatures() (*common.SeqExampleBuff, error) {
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

func (d *DeepFM) getItemExamplesFeatures() (*[]common.SeqExampleBuff, error) {
	redisKeyPrefix := d.getServiceConfig().GetModelClient().GetItemRedisKeyPre()
	itemSeqExampleBuffs := make([]common.SeqExampleBuff, 0)
	for _, itemId := range d.getItemList() {
		redisKey := redisKeyPrefix + itemId
		userExampleFeats, err := d.getServiceConfig().GetRedisClient().GetRedisPool().Get(redisKey)
		itemExampleFeatsBuff := make([]byte, 0)
		if err != nil {
			return &itemSeqExampleBuffs, nil
		} else {
			itemExampleFeatsBuff = []byte(userExampleFeats)
		}

		itemSeqExampleBuff := common.SeqExampleBuff{
			Key:  &itemId,
			Buff: &itemExampleFeatsBuff,
		}
		itemSeqExampleBuffs = append(itemSeqExampleBuffs, itemSeqExampleBuff)
	}

	return &itemSeqExampleBuffs, nil
}
