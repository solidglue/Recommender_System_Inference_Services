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

type ConfigFactory struct {
}

// faiss config factory
func (f *ConfigFactory) createFaissConfig(indexConfStr string) *faiss_config_loader.FaissIndexConfig {
	faissConfig := new(faiss_config_loader.FaissIndexConfig)
	faissConfig.ConfigLoad("", "", indexConfStr)

	return faissConfig
}

// model config factory
func (m *ConfigFactory) createModelConfig(modelConfStr string) *model_config_loader.ModelConfig {
	modelConfig := new(model_config_loader.ModelConfig)
	modelConfig.ConfigLoad("", "", modelConfStr)

	return modelConfig
}

// redis config factory
func (r *ConfigFactory) createRedisConfig(domain string, dataId string, redisConfStr string) *redis_config_loader.RedisConfig {
	redisConfig := new(redis_config_loader.RedisConfig)
	redisConfig.ConfigLoad(domain, dataId, redisConfStr)

	return redisConfig
}
