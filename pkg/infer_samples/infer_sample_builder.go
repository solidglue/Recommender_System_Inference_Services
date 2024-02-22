package infer_samples

import (
	"infer-microservices/internal/logs"
	"infer-microservices/pkg/config_loader/model_config"
	"infer-microservices/pkg/infer_samples/feature"
	"time"
)

type InferSampleBuilder struct {
	inferSample InferSample
}

// inferSample
func (b *InferSampleBuilder) SetInferSample(inferSample InferSample) {
	b.inferSample = inferSample
}

func (b *InferSampleBuilder) GetInferSample() InferSample {
	return b.inferSample
}

// TODO：特征过滤
// user offline feature build
func (b *InferSampleBuilder) UserOfflineSampleBuilder(model model_config.ModelConfig, userId string, featureList []string) *InferSampleBuilder {
	//get user offline example
	userOfflineExampleCh := make(chan *feature.SeqExampleBuff, 1)
	userExampleFeatures := &feature.SeqExampleBuff{}
	b.getUserExampleFeaturesOffline(model, userId, featureList, userOfflineExampleCh)

loop:
	for {
		select {
		case userExampleFeatures_ := <-userOfflineExampleCh:
			userExampleFeatures = userExampleFeatures_
		case <-time.After(time.Millisecond * 100):
			break loop
		}
	}

	b.inferSample.SetUserOfflineSample(userExampleFeatures)

	return b
}

// user offline feature build
func (b *InferSampleBuilder) UserRealtimeSampleBuilder(model model_config.ModelConfig, userId string, featureList []string) *InferSampleBuilder {
	//get user offline example
	userRealtimeExampleCh := make(chan *feature.SeqExampleBuff, 1)
	userExampleFeatures := &feature.SeqExampleBuff{}
	b.getUserExampleFeaturesOffline(model, userId, featureList, userRealtimeExampleCh)

loop:
	for {
		select {
		case userExampleFeatures_ := <-userRealtimeExampleCh:
			userExampleFeatures = userExampleFeatures_
		case <-time.After(time.Millisecond * 100):
			break loop
		}
	}

	b.inferSample.SetUserRealtimeSample(userExampleFeatures)

	return b
}

// user offline feature build
func (b *InferSampleBuilder) ItemsSampleBuilder(model model_config.ModelConfig, itemIdList []string, featureList []string) *InferSampleBuilder {
	//get user offline example
	itemListExampleCh := make(chan *[]feature.SeqExampleBuff, 1)
	itemExampleFeaturesList := make([]feature.SeqExampleBuff, 0)
	b.getItemExamplesFeatures(model, itemIdList, featureList, itemListExampleCh)

loop:
	for {
		select {
		case itemsExampleFeatures_ := <-itemListExampleCh:
			itemExampleFeaturesList = *itemsExampleFeatures_
		case <-time.After(time.Millisecond * 100):
			break loop
		}
	}

	b.inferSample.SetItemsSample(&itemExampleFeaturesList)

	return b
}

func (b *InferSampleBuilder) getItemExamplesFeatures(model model_config.ModelConfig, itemList []string, featureList []string, ch chan<- *[]feature.SeqExampleBuff) {
	//TODO: use bloom filter check items, avoid all items search redis.
	redisKeyPrefix := model.GetItemRedisKeyPre()
	itemSeqExampleBuffs := make([]feature.SeqExampleBuff, 0)
	itemsCh := make(chan feature.SeqExampleBuff, 100)

	for _, itemId := range itemList {
		go func(itemId string) {
			redisKey := redisKeyPrefix + itemId
			if b.inferSample.GetItemBloomFilter().Test([]byte(itemId)) {
				userExampleFeats, err := b.inferSample.serviceConfig.GetRedisConfig().GetRedisPool().Get(redisKey)
				itemExampleFeatsBuff := make([]byte, 0)
				if err != nil {
					logs.Error(err)
				} else {
					itemExampleFeatsBuff = []byte(userExampleFeats)
				}

				itemSeqExampleBuff := feature.SeqExampleBuff{
					Key:  &itemId,
					Buff: &itemExampleFeatsBuff,
				}
				itemsCh <- itemSeqExampleBuff
			}
		}(itemId)

	loop:
		for {
			select {
			case <-time.After(time.Millisecond * 100):
				break loop
			case itemCh := <-itemsCh:
				itemSeqExampleBuff := itemCh
				itemSeqExampleBuffs = append(itemSeqExampleBuffs, itemSeqExampleBuff)
			}
		}
		close(itemsCh)

	}

	ch <- &itemSeqExampleBuffs
}

// get user tfrecords offline samples
func (b *InferSampleBuilder) getUserExampleFeaturesOffline(model model_config.ModelConfig, userId string, featureList []string, ch chan<- *feature.SeqExampleBuff) {
	//INFO: use bloom filter check users, avoid all users search redis.

	userSeqExampleBuff := feature.SeqExampleBuff{}
	userExampleFeatsBuff := make([]byte, 0)

	redisKey := model.GetUserRedisKeyPreOffline() + userId
	if b.inferSample.userBloomFilter.Test([]byte(userId)) {
		userExampleFeats, err := b.inferSample.serviceConfig.GetRedisConfig().GetRedisPool().Get(redisKey)
		if err != nil {
			logs.Error("get item features err", err)
		} else {
			userExampleFeatsBuff = []byte(userExampleFeats) //.(string)
		}
	}

	//protrait features & realtime features.
	userSeqExampleBuff = feature.SeqExampleBuff{
		Key:  &userId,
		Buff: &userExampleFeatsBuff,
	}

	ch <- &userSeqExampleBuff
}

// get user tfrecords online samples
func (b *InferSampleBuilder) getUserExampleFeaturesRealtime(model model_config.ModelConfig, userId string, featureList []string, ch chan<- *feature.SeqExampleBuff) {
	//TODO: use bloom filter check users, avoid all users search redis.
	userContextSeqExampleBuff := feature.SeqExampleBuff{}
	userContextExampleFeatsBuff := make([]byte, 0)

	redisKey := model.GetUserRedisKeyPreRealtime() + userId
	if b.inferSample.userBloomFilter.Test([]byte(userId)) {
		userContextSeqExampleBuff, err := b.inferSample.serviceConfig.GetRedisConfig().GetRedisPool().Get(redisKey)
		if err != nil {
			logs.Error("get item features err", err)
		} else {
			userContextExampleFeatsBuff = []byte(userContextSeqExampleBuff) //.(string)
		}
	}

	//TODO: update context features. only from requst. such as location , time
	//context features.
	userContextSeqExampleBuff = feature.SeqExampleBuff{
		Key:  &userId,
		Buff: &userContextExampleFeatsBuff,
	}

	ch <- &userContextSeqExampleBuff
}
