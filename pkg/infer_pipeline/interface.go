package infer_pipeline

import (
	"infer-microservices/pkg/config_loader/model_config"
	feature "infer-microservices/pkg/infer_features"
	"net/http"
)

type InferPipelineInterface interface {
	//sample
	RecallSampleDirector(userId string, offlineFeature bool, onlineFeature bool, featureList []string) (feature.ExampleFeatures, error)
	RankingSampleDirector(userId string, offlineFeature bool, onlineFeature bool, itemIdList []string, featureList []string) (feature.ExampleFeatures, error)

	//infer
	ModelInferSkywalking(model model_config.ModelConfig, requestId string, exposureList []string, r *http.Request, inferSample feature.ExampleFeatures, retNum int) (map[string][]map[string]interface{}, error)
	ModelInferNoSkywalking(model model_config.ModelConfig, requestId string, exposureList []string, inferSample feature.ExampleFeatures, retNum int) (map[string][]map[string]interface{}, error)
}
