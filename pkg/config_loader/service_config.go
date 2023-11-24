package service_config_loader

import (
	"infer-microservices/pkg/config_loader/faiss_config"
	"infer-microservices/pkg/config_loader/model_config"
	"infer-microservices/pkg/config_loader/redis_config"
)

var serviceConfigs = make(map[string]*ServiceConfig, 0) //one server/dataid,one service conn

func SetServiceConfigs(serviceConfigs_ map[string]*ServiceConfig) {
	serviceConfigs = serviceConfigs_
}

func GetServiceConfigs() map[string]*ServiceConfig {
	return serviceConfigs
}

type ServiceConfig struct {
	serviceId         string                         `validate:"required,unique,min=4,max=10"` //dataid
	redisConfig       redis_config.RedisConfig       `validate:"required"`                     //redis conn info
	faissIndexConfigs faiss_config.FaissIndexConfigs //index conn info
	modelConfig       model_config.ModelConfig       `validate:"required"` //model conn info
}

func init() {
}

// serviceId
func (s *ServiceConfig) setServiceId(dataId string) {
	s.serviceId = dataId
}

func (s *ServiceConfig) GetServiceId() string {
	return s.serviceId
}

// redisConfig
func (s *ServiceConfig) setRedisConfig(redisConfig redis_config.RedisConfig) {
	s.redisConfig = redisConfig
}

func (s *ServiceConfig) GetRedisConfig() *redis_config.RedisConfig {
	return &s.redisConfig
}

// faissIndexConfig
func (s *ServiceConfig) SetFaissIndexConfigs(faissIndexConfigs faiss_config.FaissIndexConfigs) {
	s.faissIndexConfigs = faissIndexConfigs
}

func (s *ServiceConfig) GetFaissIndexConfigs() *faiss_config.FaissIndexConfigs {
	return &s.faissIndexConfigs
}

// modelConfig
func (s *ServiceConfig) setModelConfig(modelConfig model_config.ModelConfig) {
	s.modelConfig = modelConfig
}

func (s *ServiceConfig) GetModelConfig() *model_config.ModelConfig {
	return &s.modelConfig
}
