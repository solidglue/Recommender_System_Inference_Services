package infer_samples

import (
	"infer-microservices/internal/flags"
	"infer-microservices/internal/logs"
	config_loader "infer-microservices/pkg/config_loader"
	feature "infer-microservices/pkg/infer_features"

	"time"

	"github.com/allegro/bigcache"
	bloomv3 "github.com/bits-and-blooms/bloom/v3"
)


//INFO:solution-A : All tfrecored format samples have been preprocessed and stored in Redis


var bigCacheConfSample bigcache.Config
var lifeWindowS time.Duration
var cleanWindowS time.Duration
var hardMaxCacheSize int
var maxEntrySize int
var maxEntriesInWindow int
var verbose bool
var shards int
var bigCacheSample *bigcache.BigCache
var err error

type InferSample struct {
	modelName          string
	serviceConfig      *config_loader.ServiceConfig
	bigCacheSample     *bigcache.BigCache
	userBloomFilter    *bloomv3.BloomFilter
	itemBloomFilter    *bloomv3.BloomFilter
	userOfflineSample  *feature.SeqExampleBuff
	userRealtimeSample *feature.SeqExampleBuff
	itemsSample        *[]feature.SeqExampleBuff
}

func init() {
	flagFactory := flags.FlagFactory{}
	flagCache := flagFactory.CreateFlagCache()
	lifeWindowS = time.Duration(*flagCache.GetBigcacheLifeWindowS())
	cleanWindowS = time.Duration(*flagCache.GetBigcacheCleanWindowS())
	hardMaxCacheSize = *flagCache.GetBigcacheHardMaxCacheSize()
	maxEntrySize = *flagCache.GetBigcacheMaxEntrySize()
	bigCacheConfSample = bigcache.Config{
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

	//set cache
	bigCacheSample, err = bigcache.NewBigCache(bigCacheConfSample)
	if err != nil {
		logs.Error(err)
	}
}

// userid
func (b InferSample) SetModelName(modelName string) {
	b.modelName = modelName
}

func (b InferSample) GetModelName() string {
	return b.modelName
}

// serviceConfig *service_config.ServiceConfig
func (b InferSample) SetServiceConfig(serviceConfig *config_loader.ServiceConfig) {
	b.serviceConfig = serviceConfig
}

func (b InferSample) GetServiceConfig() *config_loader.ServiceConfig {
	return b.serviceConfig
}

func (s InferSample) SetUserBloomFilter(filter *bloomv3.BloomFilter) {
	s.userBloomFilter = filter
}

func (s InferSample) GetUserBloomFilter() *bloomv3.BloomFilter {
	return s.userBloomFilter
}

func (s InferSample) SetItemBloomFilter(filter *bloomv3.BloomFilter) {
	s.itemBloomFilter = filter
}

func (s InferSample) GetItemBloomFilter() *bloomv3.BloomFilter {
	return s.itemBloomFilter
}

func (s InferSample) GetUserOfflineSample() *feature.SeqExampleBuff {
	return s.userOfflineSample
}

func (s InferSample) SetUserOfflineSample(userOfflineSample *feature.SeqExampleBuff) {
	s.userOfflineSample = userOfflineSample
}

func (s InferSample) GetUserRealtimeSample() *feature.SeqExampleBuff {
	return s.userRealtimeSample
}

func (s InferSample) SetUserRealtimeSample(userRealtimeSample *feature.SeqExampleBuff) {
	s.userRealtimeSample = userRealtimeSample
}

func (s InferSample) GetItemsSample() *[]feature.SeqExampleBuff {
	return s.itemsSample
}

func (s InferSample) SetItemsSample(itemsSample *[]feature.SeqExampleBuff) {
	s.itemsSample = itemsSample
}

// observer nontify
func (s InferSample) notify(sub Subject) {
	//reload baseModel
	s.SetUserBloomFilter(GetUserBloomFilterInstance())
	s.SetItemBloomFilter(GetItemBloomFilterInstance())
}
