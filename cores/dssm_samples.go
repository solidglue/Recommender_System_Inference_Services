package cores

import (
	"encoding/json"
	"infer-microservices/common"
	"infer-microservices/common/flags"
	"infer-microservices/utils"
	"infer-microservices/utils/logs"
	"time"

	//example "infer-microservices/common/tensorflow_gogofaster/example"
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

	// lifeWindowS = time.Duration(*flags.Bigcache_lifeWindowS)
	// cleanWindowS = time.Duration(*flags.Bigcache_cleanWindowS)
	// hardMaxCacheSize = *flags.Bigcache_hardMaxCacheSize
	// maxEntrySize = *flags.Bigcache_maxEntrySize

	flagFactory := flags.FlagFactory{}
	flagCache := flagFactory.FlagCacheFactory()

	lifeWindowS = time.Duration(*flagCache.GetBigcacheLifeWindowS())
	cleanWindowS = time.Duration(*flagCache.GetBigcacheCleanWindowS())
	hardMaxCacheSize = *flagCache.GetBigcacheHardMaxCacheSize()
	maxEntrySize = *flagCache.GetBigcacheMaxEntrySize()

	//how to config : http://liuqh.icu/2021/06/15/go/package/14-bigcache/

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

func (d *Dssm) getInferExampleFeatures() (common.ExampleFeatures, error) {

	// //TODO:如果不直接查询redis取特征，需要用到feature.pb.go去构造tfrecord的example。
	// //对于从redis取特征的有2种方案，单独一个特征进程实时构造样本，然后写入redis。go服务只使用。方案2：离线特征redis查询，实时特征go构造
	// //倾向于方案1，因为实时样本可以存档，实时样本堆积就是离线样本，这样就不用每天再跑全量的离线样本了（不跑或跑部分）。且可以降低推理时间

	// 	//Feature_FloatList结构体，封装了FloatList
	// //   int64_list=tf.train.Int64List(value=[28])
	// //此处用的接口 ，Feature_FloatList等实现了isFeature_Kind接口
	// //   tf.train.Feature(int64_list=tf.train.Int64List(value=[28])),
	// //  "height": tf.train.Feature(int64_list=tf.train.Int64List(value=[28])),

	// f := make(map[string]*example.Feature, 0)
	// f[field_name+"_id"] = &example.Feature{Kind: &example.Feature_Int64List{Int64List: &example.Int64List{Value: int64_values}}}  //多值特征的话，int64_values 是特征数组
	// f[field_name+"_value"] = &example.Feature{Kind: &example.Feature_FloatList{FloatList: &example.FloatList{Value: float_values}}}  //多值特征的话，float_values是特征权重数组

	// example_f := &example.Features{
	// 	Feature: f,
	// }

	//var cacheTimeSecond = time.Duration(featuresCacheTimeInt)
	//var goCache = cache.New(cacheTimeSecond*time.Second, 2*cacheTimeSecond*time.Second)
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

// TODO: 接口实现样本构造，如果不，则把DSSM类拆分长推理和样本2部分。样本和推理写在一个类的违反单一原则

func (d *Dssm) getUserExampleFeatures() (*common.SeqExampleBuff, error) {

	//TODO: update context features.

	redisKey := d.getServiceConfig().GetModelClient().GetUserRedisKeyPre() + d.getUserId()
	userExampleFeats, err := d.getServiceConfig().GetRedisClient().GetRedisPool().Get(redisKey)

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

// TODO:此处改成多态，减少代码冗余
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
