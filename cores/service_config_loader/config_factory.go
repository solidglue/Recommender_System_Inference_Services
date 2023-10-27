package service_config_loader

import (
	"infer-microservices/cores/service_config_loader/faiss_config_loader"
	"infer-microservices/cores/service_config_loader/model_config_loader"
	"infer-microservices/cores/service_config_loader/redis_config_loader"
)

// config interface
type ConfigLoadInterface interface {
	ConfigLoad(domain string, dataId string, confStr string) error
}

type faissConfigFactory struct {
}

type modelConfigFactory struct {
}

type redisConfigFactory struct {
}

// faiss config factory
func (f *faissConfigFactory) createConfigLoader(indexConfStr string) *faiss_config_loader.FaissIndexConfig {
	faissLoader := new(faiss_config_loader.FaissIndexConfig)
	faissLoader.ConfigLoad("", "", indexConfStr)

	return faissLoader
}

// model config factory
func (m *modelConfigFactory) createConfigLoader(modelConfStr string) *model_config_loader.ModelConfig {
	modelLoader := new(model_config_loader.ModelConfig)
	modelLoader.ConfigLoad("", "", modelConfStr)

	return modelLoader
}

// redis config factory
func (r *redisConfigFactory) createConfigLoader(domain string, dataId string, redisConfStr string) *redis_config_loader.RedisConfig {
	redisLoader := new(redis_config_loader.RedisConfig)
	redisLoader.ConfigLoad(domain, dataId, redisConfStr)

	return redisLoader
}
