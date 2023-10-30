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

func (b *ServiceConfigBuilder) RedisConfigBuilder(domain string, DataId string, redisConfStr string) *ServiceConfigBuilder {
	configFactory := &ConfigFactory{}
	redisConfig := configFactory.createRedisConfig(domain, DataId, redisConfStr)
	b.serviceConfig.setRedisConfig(*redisConfig)

	return b
}

func (b *ServiceConfigBuilder) FaissConfigBuilder(domain string, DataId string, indexConfStr string) *ServiceConfigBuilder {
	//load redis conf
	configFactory := &ConfigFactory{}
	faissConfig := configFactory.createFaissConfig(indexConfStr)
	b.serviceConfig.setFaissIndexConfig(*faissConfig)

	return b
}

func (b *ServiceConfigBuilder) ModelConfigBuilder(domain string, DataId string, modelConfStr string) *ServiceConfigBuilder {
	//load redis conf
	configFactory := &ConfigFactory{}
	modelConfig := configFactory.createModelConfig(modelConfStr)
	b.serviceConfig.setModelConfig(*modelConfig)

	return b
}
