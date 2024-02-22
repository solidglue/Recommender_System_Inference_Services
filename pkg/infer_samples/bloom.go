package infer_samples

import (
	"infer-microservices/internal/flags"

	bloomv3 "github.com/bits-and-blooms/bloom/v3"
)

var userBloomFilterInstance *bloomv3.BloomFilter
var itemBloomFilterInstance *bloomv3.BloomFilter

func init() {
	flagFactory := flags.FlagFactory{}
	flagBloom := flagFactory.CreateFlagBloom()
	userBloomFilterInstance = bloomv3.NewWithEstimates(*flagBloom.GetUserCountLevel(), 0.01)
	itemBloomFilterInstance = bloomv3.NewWithEstimates(*flagBloom.GetItemCountLevel(), 0.01)
}

// INFO: singleton instance
func GetUserBloomFilterInstance() *bloomv3.BloomFilter {
	return userBloomFilterInstance
}

// INFO: singleton instance
func GetItemBloomFilterInstance() *bloomv3.BloomFilter {
	return itemBloomFilterInstance
}

func BloomPush(filter *bloomv3.BloomFilter, id string) {
	filter.Add([]byte(id))
}

func BloomClean(filter *bloomv3.BloomFilter) {
	filter.ClearAll()
}
