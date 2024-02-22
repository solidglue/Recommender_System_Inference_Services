package infer_model

import (
	"infer-microservices/pkg/config_loader/model_config"
	"infer-microservices/pkg/infer_samples/feature"

	"net/http"
)

type ModelStrategyContext struct {
	modelStrategy ModelStrategyInterface
}

func (m *ModelStrategyContext) SetModelStrategy(strategy ModelStrategyInterface) {
	m.modelStrategy = strategy
}

func (m *ModelStrategyContext) ModelInferSkywalking(model model_config.ModelConfig, requestId string, r *http.Request, inferSample feature.ExampleFeatures, retNum int) (map[string][]map[string]interface{}, error) {
	response, err := m.modelStrategy.ModelInferSkywalking(model, requestId, r, inferSample, retNum)
	return response, err
}

func (m *ModelStrategyContext) ModelInferNoSkywalking(model model_config.ModelConfig, requestId string, inferSample feature.ExampleFeatures, retNum int) (map[string][]map[string]interface{}, error) {
	response, err := m.modelStrategy.ModelInferNoSkywalking(model, requestId, inferSample, retNum)
	return response, err
}
