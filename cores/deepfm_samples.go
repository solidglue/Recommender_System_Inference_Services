package cores

import (
	"encoding/json"
	"infer-microservices/common"
	"infer-microservices/common/flags"

	//example "infer-microservices/common/tensorflow_gogofaster/example"
	"infer-microservices/utils"
	"time"

	"github.com/allegro/bigcache"
)

//TODO:用接口合并deepfm和dssm样本生成.特征类？

//TODO:本地缓存可能有问题，第一次请求缓存后，如果第二次请求负载到其它机器，则缓存失效。分布式缓存？或者redis
//理论上第二次请求用户行为已改变，用不上缓存，也没问题

var lifeWindowS1 time.Duration
var cleanWindowS1 time.Duration
var hardMaxCacheSize1 int
var maxEntrySize1 int
var maxEntriesInWindow1 int
var verbose1 bool
var shards1 int

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

func (d *DeepFM) getInferExampleFeatures() (common.ExampleFeatures, error) {

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
	if lifeWindowS > 0 {
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

	if lifeWindowS > 0 {
		bigCache.Set(cacheKeyPrefix, []byte(utils.Struct2Json(exampleData)))
	}

	return exampleData, nil
}

// TODO: 接口实现样本构造，如果不，则把DEEMPFM类拆分长推理和样本2部分。样本和推理写在一个类的违反单一原则
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
