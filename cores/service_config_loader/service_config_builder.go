package service_config_loader

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

// redis builder
func (b *ServiceConfigBuilder) RedisConfigBuilder(dataId string, redisConfStr string) *ServiceConfigBuilder {
	configFactory := &ConfigFactory{}
	redisConfig := configFactory.createRedisConfig(dataId, redisConfStr)
	b.serviceConfig.setRedisConfig(*redisConfig)

	return b
}

// faiss builder
func (b *ServiceConfigBuilder) FaissConfigBuilder(dataId string, indexConfStr string) *ServiceConfigBuilder {
	configFactory := &ConfigFactory{}
	faissConfig := configFactory.createFaissConfig(dataId, indexConfStr)
	b.serviceConfig.setFaissIndexConfig(*faissConfig)

	return b
}

// model builder
func (b *ServiceConfigBuilder) ModelConfigBuilder(domain string, dataId string, modelConfStr string) *ServiceConfigBuilder {
	configFactory := &ConfigFactory{}
	modelConfig := configFactory.createModelConfig(dataId, modelConfStr)
	b.serviceConfig.setModelConfig(*modelConfig)

	return b
}
