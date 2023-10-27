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
	configFactory := &redisConfigFactory{}
	redisConfig := configFactory.createConfigLoader(domain, DataId, redisConfStr)
	b.serviceConfig.setRedisConfig(*redisConfig)

	return b
}

func (b *ServiceConfigBuilder) FaissConfigBuilder(domain string, DataId string, indexConfStr string) *ServiceConfigBuilder {
	//load redis conf
	configFactory := &faissConfigFactory{}
	faissConfig := configFactory.createConfigLoader(indexConfStr)
	b.serviceConfig.setFaissIndexConfig(*faissConfig)

	return b
}

func (b *ServiceConfigBuilder) ModelConfigBuilder(domain string, DataId string, modelConfStr string) *ServiceConfigBuilder {
	//load redis conf
	configFactory := &modelConfigFactory{}
	modelConfig := configFactory.createConfigLoader(modelConfStr)
	b.serviceConfig.setModelConfig(*modelConfig)

	return b
}
