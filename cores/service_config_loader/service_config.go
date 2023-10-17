package service_config_loader

import (
	"infer-microservices/common/flags"
	"infer-microservices/cores/service_config_loader/faiss_config_loader"
	"infer-microservices/cores/service_config_loader/model_config_loader"
	"infer-microservices/cores/service_config_loader/redis_config_loader"
)

var serviceConfFile string
var serviceConfigInstance *ServiceConfig

type ServiceConfig struct {
	serviceId        string                               //dataid
	redisClient      redis_config_loader.RedisClient      //redis conn info
	faissIndexClient faiss_config_loader.FaissIndexClient //index conn info
	modelClient      model_config_loader.ModelClient      //model conn info
}

func init() {
	flagFactory := flags.FlagFactory{}
	flagServiceConfig := flagFactory.FlagServiceConfigFactory()
	serviceConfFile = *flagServiceConfig.GetServiceConfigFile()
}

// serviceId
func (s *ServiceConfig) setServiceId(dataId string) {
	s.serviceId = dataId
}

func (s *ServiceConfig) GetServiceId() string {
	return s.serviceId
}

// redisClient
func (s *ServiceConfig) setRedisClient(redisClient redis_config_loader.RedisClient) {
	s.redisClient = redisClient
}

func (s *ServiceConfig) GetRedisClient() *redis_config_loader.RedisClient {
	return &s.redisClient
}

// faissIndexClient
func (s *ServiceConfig) setFaissIndexClient(faissIndexClient faiss_config_loader.FaissIndexClient) {
	s.faissIndexClient = faissIndexClient
}

func (s *ServiceConfig) GetFaissIndexClient() *faiss_config_loader.FaissIndexClient {
	return &s.faissIndexClient
}

// modelClient
func (s *ServiceConfig) setModelClient(modelClient model_config_loader.ModelClient) {
	s.modelClient = modelClient
}

func (s *ServiceConfig) GetModelClient() *model_config_loader.ModelClient {
	return &s.modelClient
}
