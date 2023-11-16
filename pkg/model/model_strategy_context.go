package model

import "net/http"

type ModelStrategyContext struct {
	modelStrategy ModelStrategyInterface
}

func (m *ModelStrategyContext) SetModelStrategy(strategy ModelStrategyInterface) {
	m.modelStrategy = strategy
}

func (m *ModelStrategyContext) ModelInferSkywalking(r *http.Request) (map[string]interface{}, error) {
	response, err := m.modelStrategy.ModelInferSkywalking(r)
	return response, err
}

func (m *ModelStrategyContext) ModelInferNoSkywalking(r *http.Request) (map[string]interface{}, error) {
	response, err := m.modelStrategy.ModelInferNoSkywalking(r)
	return response, err
}
