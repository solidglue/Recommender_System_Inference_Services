package service_config_loader

type ServiceConfigDirector struct {
}

// build contain index
func (s *ServiceConfigDirector) ServiceConfigUpdateContainIndexDirector(domain string, dataId string,
	redisConfStr string, modelConfStr string, indexConfStr string) ServiceConfig {
	serviceConfigBuilder := ServiceConfigBuilder{}
	builder := serviceConfigBuilder.RedisConfigBuilder(dataId, redisConfStr).FaissConfigBuilder(dataId, indexConfStr).ModelConfigBuilder(domain, dataId, modelConfStr)
	serviceConfig := builder.GetServiceConfig()
	serviceConfig.setServiceId(dataId)

	return serviceConfig
}

// build not contain index
func (s *ServiceConfigDirector) ServiceConfigUpdaterNotContainIndexDirector(domain string, dataId string,
	redisConfStr string, modelConfStr string) ServiceConfig {
	serviceConfigBuilder := ServiceConfigBuilder{}
	builder := serviceConfigBuilder.RedisConfigBuilder(dataId, redisConfStr).ModelConfigBuilder(domain, dataId, modelConfStr)
	serviceConfig := builder.GetServiceConfig()
	serviceConfig.setServiceId(dataId)

	return serviceConfig
}
