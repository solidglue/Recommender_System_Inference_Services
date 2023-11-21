package service_config_loader

import (
	"infer-microservices/pkg/config_loader/faiss_config"
	"infer-microservices/pkg/config_loader/model_config"
	"infer-microservices/pkg/config_loader/redis_config"
)

// config interface
type ConfigLoadInterface interface {
	ConfigLoad(domain string, dataId string, confStr string) error
}

type ConfigFactory struct {
}

// faiss config factory
func (f *ConfigFactory) createFaissConfig(dataId string, indexConfStr string) *faiss_config.FaissIndexConfigs {
	faissConfigs := new(faiss_config.FaissIndexConfigs)
	faissConfigs.ConfigLoad(dataId, indexConfStr)

	return faissConfigs
}

// model config factory
func (m *ConfigFactory) createModelConfig(dataId string, modelConfStr string) *model_config.ModelConfig {
	modelConfig := new(model_config.ModelConfig)
	modelConfig.ConfigLoad(dataId, modelConfStr)

	return modelConfig
}

// redis config factory
func (r *ConfigFactory) createRedisConfig(dataId string, redisConfStr string) *redis_config.RedisConfig {
	redisConfig := new(redis_config.RedisConfig)
	redisConfig.ConfigLoad(dataId, redisConfStr)

	return redisConfig
}
