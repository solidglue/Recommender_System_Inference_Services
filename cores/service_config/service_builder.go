package service_config

import (
	"infer-microservices/cores/faiss"
	"infer-microservices/cores/model"
	"infer-microservices/cores/redis_config"
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
	redisFactory := redis_config.RedisFactory{}
	redisClient := redisFactory.RedisClientFactory(domain, DataId, redisConfStr)
	b.serviceConfig.setRedisClient(*redisClient)

	return b
}

func (b *ServiceConfigBuilder) FaissClientBuilder(indexConfStr string) *ServiceConfigBuilder {
	//load redis conf
	faissFactory := faiss.FaissFactory{}
	faissClient := faissFactory.FaissClientFactory(indexConfStr)
	b.serviceConfig.setFaissIndexClient(*faissClient)

	return b
}

func (b *ServiceConfigBuilder) ModelClientBuilder(modelConfStr string) *ServiceConfigBuilder {
	//load redis conf
	modelFactory := model.ModelFactory{}
	modelClient := modelFactory.ModelClientFactory(modelConfStr)
	b.serviceConfig.setModelClient(*modelClient)

	return b
}
