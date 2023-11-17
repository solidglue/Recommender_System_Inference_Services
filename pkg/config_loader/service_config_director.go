package service_config_loader

import (
	"infer-microservices/pkg/logs"
	"time"

	validator "github.com/go-playground/validator/v10"
)

type ServiceConfigDirector struct {
	configBuilder ServiceConfigBuilder
}

// set
func (s *ServiceConfigDirector) SetConfigBuilder(builder ServiceConfigBuilder) {
	s.configBuilder = builder
}

// build contain index
func (s *ServiceConfigDirector) ServiceConfigUpdateContainIndexDirector(dataId string,
	redisConfStr string, modelConfStr string, indexConfStr string) ServiceConfig {
	builder := s.configBuilder.RedisConfigBuilder(dataId, redisConfStr).FaissConfigBuilder(dataId, indexConfStr).ModelConfigBuilder(dataId, modelConfStr)
	serviceConfig := builder.GetServiceConfig()
	serviceConfig.setServiceId(dataId)

	//validete serviceConfig
	validate := validator.New()
	err := validate.Struct(serviceConfig)
	if err != nil {
		logs.Error(dataId, time.Now(), err)
		return ServiceConfig{}
	}

	return serviceConfig
}

// build not contain index
func (s *ServiceConfigDirector) ServiceConfigUpdaterNotContainIndexDirector(dataId string,
	redisConfStr string, modelConfStr string) ServiceConfig {
	builder := s.configBuilder.RedisConfigBuilder(dataId, redisConfStr).ModelConfigBuilder(dataId, modelConfStr)
	serviceConfig := builder.GetServiceConfig()
	serviceConfig.setServiceId(dataId)

	//validete serviceConfig
	validate := validator.New()
	err := validate.Struct(serviceConfig)
	if err != nil {
		logs.Error(dataId, time.Now(), err)
		return ServiceConfig{}
	}

	return serviceConfig
}
