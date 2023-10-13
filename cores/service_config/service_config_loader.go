package service_config

import (
	"infer-microservices/common/flags"
	"infer-microservices/cores/faiss"
	"infer-microservices/cores/model"
	"infer-microservices/cores/service_config/redis_config"
)

var serviceConfFile string
var serviceConfigInstance *ServiceConfig

type ServiceConfig struct {
	serviceId        string                   //dataid
	redisClient      redis_config.RedisClient //redis conn info
	faissIndexClient faiss.FaissIndexClient   //index conn info
	modelClient      model.ModelClient        //model conn info
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
func (s *ServiceConfig) setRedisClient(redisClient redis_config.RedisClient) {
	s.redisClient = redisClient
}

func (s *ServiceConfig) GetRedisClient() *redis_config.RedisClient {
	return &s.redisClient
}

// faissIndexClient
func (s *ServiceConfig) setFaissIndexClient(faissIndexClient faiss.FaissIndexClient) {
	s.faissIndexClient = faissIndexClient
}

func (s *ServiceConfig) GetFaissIndexClient() *faiss.FaissIndexClient {
	return &s.faissIndexClient
}

// modelClient
func (s *ServiceConfig) setModelClient(modelClient model.ModelClient) {
	s.modelClient = modelClient
}

func (s *ServiceConfig) GetModelClient() *model.ModelClient {
	return &s.modelClient
}
