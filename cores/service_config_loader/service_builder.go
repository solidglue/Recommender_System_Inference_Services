package service_config_loader

import (
	"infer-microservices/cores/service_config_loader/faiss_config_loader"
	"infer-microservices/cores/service_config_loader/model_config_loader"
	"infer-microservices/cores/service_config_loader/redis_config_loader"
)

type ServiceConfigBuilder struct {
	serviceConfig ServiceConfig
}

// serviceConfig
func (b *ServiceConfigBuilder) SetServiceConfig(serviceConfig ServiceConfig) {
	b.serviceConfig = serviceConfig
}

func (b *ServiceConfigBuilder) GetServiceConfig() ServiceConfig {
	return b.serviceConfig
}

func (b *ServiceConfigBuilder) RedisClientBuilder(domain string, DataId string, redisConfStr string) *ServiceConfigBuilder {
	redisFactory := redis_config_loader.RedisFactory{}
	redisClient := redisFactory.RedisClientFactory(domain, DataId, redisConfStr)
	b.serviceConfig.setRedisClient(*redisClient)

	return b
}

func (b *ServiceConfigBuilder) FaissClientBuilder(indexConfStr string) *ServiceConfigBuilder {
	//load redis conf
	faissFactory := faiss_config_loader.FaissFactory{}
	faissClient := faissFactory.FaissClientFactory(indexConfStr)
	b.serviceConfig.setFaissIndexClient(*faissClient)

	return b
}

func (b *ServiceConfigBuilder) ModelClientBuilder(modelConfStr string) *ServiceConfigBuilder {
	//load redis conf
	modelFactory := model_config_loader.ModelFactory{}
	modelClient := modelFactory.ModelClientFactory(modelConfStr)
	b.serviceConfig.setModelClient(*modelClient)

	return b
}
