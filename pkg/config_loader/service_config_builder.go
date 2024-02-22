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
	faissConfigs := configFactory.createFaissConfig(dataId, indexConfStr)
	b.serviceConfig.SetFaissIndexConfigs(*faissConfigs)

	return b
}

// model builder
func (b *ServiceConfigBuilder) ModelConfigBuilder(dataId string, modelConfStr string) *ServiceConfigBuilder {
	configFactory := &ConfigFactory{}
	modelConfigs := configFactory.createModelConfig(dataId, modelConfStr)
	b.serviceConfig.setModelsConfig(*modelConfigs)

	return b
}

// pipeline builder
func (b *ServiceConfigBuilder) PipelineConfigBuilder(dataId string, pipelineConfStr string) *ServiceConfigBuilder {
	configFactory := &ConfigFactory{}
	PipelineConfig := configFactory.createPipelineConfig(dataId, pipelineConfStr)
	b.serviceConfig.setPipelineConfig(*PipelineConfig)

	return b
}
