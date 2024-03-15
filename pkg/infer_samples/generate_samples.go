package infer_samples

import (
	config_loader "infer-microservices/pkg/config_loader"
	feature "infer-microservices/pkg/infer_features"

	"github.com/allegro/bigcache"
	bloomv3 "github.com/bits-and-blooms/bloom/v3"
)

//TODO: ADD
//INFO: solution-B：Query 、process and build tfrecord samples during inference.

type InferSampleGenerate struct {
	modelName       string
	serviceConfig   *config_loader.ServiceConfig
	bigCacheSample  *bigcache.BigCache
	userBloomFilter *bloomv3.BloomFilter
	itemBloomFilter *bloomv3.BloomFilter
	inferFeature    feature.InferFeature // all the features using to generate tfrecord samples.
}
