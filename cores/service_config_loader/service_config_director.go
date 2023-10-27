package service_config_loader

type ServiceConfigDirector struct {
}

func (s *ServiceConfigDirector) ServiceConfigUpdateContainIndexDirector(domain string, dataId string,
	redisConfStr string, modelConfStr string, indexConfStr string) ServiceConfig {
	//load redis,faiss,model
	serviceConfigBuilder := ServiceConfigBuilder{}
	builder := serviceConfigBuilder.RedisConfigBuilder(domain, dataId, redisConfStr).FaissConfigBuilder(domain, dataId, indexConfStr).ModelConfigBuilder(domain, dataId, modelConfStr)

	return builder.GetServiceConfig()
}

func (s *ServiceConfigDirector) ServiceConfigUpdaterNotContainIndexDirector(domain string, dataId string,
	redisConfStr string, modelConfStr string) ServiceConfig {
	//load redis,model
	serviceConfigBuilder := ServiceConfigBuilder{}
	builder := serviceConfigBuilder.RedisConfigBuilder(domain, dataId, redisConfStr).ModelConfigBuilder(domain, dataId, modelConfStr)

	return builder.GetServiceConfig()
}
