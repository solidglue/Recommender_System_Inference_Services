package infer_model

import (
	"infer-microservices/pkg/config_loader/model_config"
	feature "infer-microservices/pkg/infer_features"

	"net/http"
)

type ModelStrategyContext struct {
	modelStrategy ModelStrategyInterface
}

func (m *ModelStrategyContext) SetModelStrategy(strategy ModelStrategyInterface) {
	m.modelStrategy = strategy
}

func (m *ModelStrategyContext) ModelInferSkywalking(model model_config.ModelConfig, requestId string, exposureList []string, r *http.Request, inferSample feature.ExampleFeatures, retNum int) (map[string][]map[string]interface{}, error) {
	response, err := m.modelStrategy.ModelInferSkywalking(model, requestId, exposureList, r, inferSample, retNum)
	return response, err
}

func (m *ModelStrategyContext) ModelInferNoSkywalking(model model_config.ModelConfig, requestId string, exposureList []string, inferSample feature.ExampleFeatures, retNum int) (map[string][]map[string]interface{}, error) {
	response, err := m.modelStrategy.ModelInferNoSkywalking(model, requestId, exposureList, inferSample, retNum)
	return response, err
}
