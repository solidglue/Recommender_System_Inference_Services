package service_config_loader

type ServiceConfigDirector struct {
	configBuilder ServiceConfigBuilder
}

// set
func (s *ServiceConfigDirector) SetConfigBuilder(builder ServiceConfigBuilder) {
	s.configBuilder = builder
}

// build contain index
func (s *ServiceConfigDirector) ServiceConfigUpdateContainIndexDirector(domain string, dataId string,
	redisConfStr string, modelConfStr string, indexConfStr string) ServiceConfig {
	builder := s.configBuilder.RedisConfigBuilder(dataId, redisConfStr).FaissConfigBuilder(dataId, indexConfStr).ModelConfigBuilder(domain, dataId, modelConfStr)
	serviceConfig := builder.GetServiceConfig()
	serviceConfig.setServiceId(dataId)

	return serviceConfig
}

// build not contain index
func (s *ServiceConfigDirector) ServiceConfigUpdaterNotContainIndexDirector(domain string, dataId string,
	redisConfStr string, modelConfStr string) ServiceConfig {
	builder := s.configBuilder.RedisConfigBuilder(dataId, redisConfStr).ModelConfigBuilder(domain, dataId, modelConfStr)
	serviceConfig := builder.GetServiceConfig()
	serviceConfig.setServiceId(dataId)

	return serviceConfig
}
