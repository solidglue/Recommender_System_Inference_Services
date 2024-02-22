package infer_samples

import (
	"encoding/json"
	"infer-microservices/internal/logs"
	"infer-microservices/internal/utils"
	"infer-microservices/pkg/config_loader/model_config"
	"infer-microservices/pkg/infer_samples/feature"
	"net/http"
)

type InferSampleDirector struct {
	inferSampleBuilder InferSampleBuilder
}

// set
func (s *InferSampleDirector) SetConfigBuilder(builder InferSampleBuilder) {
	s.inferSampleBuilder = builder
}

// Each model may have multiple ways to create samples, using callback functions to determine which method to call
func (s *InferSampleDirector) RecallSampleDirector(model model_config.ModelConfig, userId string, offlineFeature bool, onlineFeature bool, featureList []string) (feature.ExampleFeatures, error) {
	cacheKeyPrefix := userId + s.inferSampleBuilder.inferSample.serviceConfig.GetServiceId() + "_samples"
	userOfflineExampleFeatures := &feature.SeqExampleBuff{}
	userRealtimeExampleFeatures := &feature.SeqExampleBuff{}

	exampleData := feature.ExampleFeatures{
		UserExampleFeatures:        userOfflineExampleFeatures,
		UserContextExampleFeatures: userRealtimeExampleFeatures,
	}
	exampleDataBytes, _ := s.inferSampleBuilder.inferSample.bigCacheSample.Get(cacheKeyPrefix)
	err = json.Unmarshal(exampleDataBytes, &exampleData)
	if err != nil {
		logs.Error(err)
		return exampleData, nil

	}

	if offlineFeature && onlineFeature {
		builder := s.inferSampleBuilder.UserOfflineSampleBuilder(model, userId, featureList).UserRealtimeSampleBuilder(model, userId, featureList)
		s.inferSampleBuilder = *builder
	} else if offlineFeature && onlineFeature == false {
		builder := s.inferSampleBuilder.UserOfflineSampleBuilder(model, userId, featureList)
		s.inferSampleBuilder = *builder
	} else if offlineFeature == false && onlineFeature {
		builder := s.inferSampleBuilder.UserRealtimeSampleBuilder(model, userId, featureList)
		s.inferSampleBuilder = *builder
	}

	exampleData = feature.ExampleFeatures{
		UserExampleFeatures:        s.inferSampleBuilder.inferSample.userOfflineSample,
		UserContextExampleFeatures: s.inferSampleBuilder.inferSample.userRealtimeSample,
	}
	s.inferSampleBuilder.inferSample.bigCacheSample.Set(cacheKeyPrefix, []byte(utils.ConvertStructToJson(exampleData)))

	inferSampleInstance := s.inferSampleBuilder.inferSample
	inferSampleSubject.AddObserver(inferSampleInstance) //TODO: 构造样本的时候加入观察者，每个模型都加入

	return exampleData, nil
}

// includ pre-raning and ranking
func (s *InferSampleDirector) RankingSampleDirector(model model_config.ModelConfig, userId string, offlineFeature bool, onlineFeature bool, itemIdList []string, featureList []string) (feature.ExampleFeatures, error) {
	cacheKeyPrefix := userId + s.inferSampleBuilder.inferSample.serviceConfig.GetServiceId() + "_samples"
	userOfflineExampleFeatures := &feature.SeqExampleBuff{}
	userRealtimeExampleFeatures := &feature.SeqExampleBuff{}
	itemSeqExampleFeatures := make([]feature.SeqExampleBuff, 0)
	exampleData := feature.ExampleFeatures{
		UserExampleFeatures:        userOfflineExampleFeatures,
		UserContextExampleFeatures: userRealtimeExampleFeatures,
		ItemSeqExampleFeatures:     &itemSeqExampleFeatures,
	}
	exampleDataBytes, _ := s.inferSampleBuilder.inferSample.bigCacheSample.Get(cacheKeyPrefix)
	err = json.Unmarshal(exampleDataBytes, &exampleData)
	if err != nil {
		logs.Error(err)
		return exampleData, nil

	}

	if offlineFeature && onlineFeature {
		builder := s.inferSampleBuilder.UserOfflineSampleBuilder(model, userId, featureList).UserRealtimeSampleBuilder(model, userId, featureList).ItemsSampleBuilder(model, itemIdList, featureList)
		s.inferSampleBuilder = *builder
	} else if offlineFeature && onlineFeature == false {
		builder := s.inferSampleBuilder.UserOfflineSampleBuilder(model, userId, featureList).ItemsSampleBuilder(model, itemIdList, featureList)
		s.inferSampleBuilder = *builder
	} else if offlineFeature == false && onlineFeature {
		builder := s.inferSampleBuilder.UserRealtimeSampleBuilder(model, userId, featureList).ItemsSampleBuilder(model, itemIdList, featureList)
		s.inferSampleBuilder = *builder
	}

	exampleData = feature.ExampleFeatures{
		UserExampleFeatures:        s.inferSampleBuilder.inferSample.userOfflineSample,
		UserContextExampleFeatures: s.inferSampleBuilder.inferSample.userRealtimeSample,
		ItemSeqExampleFeatures:     s.inferSampleBuilder.inferSample.itemsSample,
	}
	s.inferSampleBuilder.inferSample.bigCacheSample.Set(cacheKeyPrefix, []byte(utils.ConvertStructToJson(exampleData)))

	return exampleData, nil
}

func (s *InferSampleDirector) ModelInferSkywalking(model model_config.ModelConfig, requestId string, r *http.Request, inferSample feature.ExampleFeatures, retNum int) (map[string][]map[string]interface{}, error) {
	//panic("no not need to Implement interface methods")
	response := make(map[string][]map[string]interface{}, 0)
	return response, nil
}

func (s *InferSampleDirector) ModelInferNoSkywalking(model model_config.ModelConfig, requestId string, inferSample feature.ExampleFeatures, retNum int) (map[string][]map[string]interface{}, error) {
	//panic("no not need to Implement interface methods")
	response := make(map[string][]map[string]interface{}, 0)
	return response, nil
}
