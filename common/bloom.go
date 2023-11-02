package common

import (
	"infer-microservices/common/flags"

	bloom "github.com/bits-and-blooms/bloom/v3"
)

var userBloomFilterInstance *bloom.BloomFilter
var itemBloomFilterInstance *bloom.BloomFilter

func init() {
	flagFactory := flags.FlagFactory{}
	flagBloom := flagFactory.CreateFlagBloom()
	userBloomFilterInstance = bloom.NewWithEstimates(*flagBloom.GetUserCountLevel(), 0.01)
	itemBloomFilterInstance = bloom.NewWithEstimates(*flagBloom.GetItemCountLevel(), 0.01)
}

// INFO: singleton instance
func GetUserBloomFilterInstance() *bloom.BloomFilter {
	return userBloomFilterInstance
}

// INFO: singleton instance
func GetItemBloomFilterInstance() *bloom.BloomFilter {
	return itemBloomFilterInstance
}

func bloomPush(filter *bloom.BloomFilter, id string) {
	filter.Add([]byte(id))
}

func CleanBloom(filter *bloom.BloomFilter) {
	filter.ClearAll()
}
