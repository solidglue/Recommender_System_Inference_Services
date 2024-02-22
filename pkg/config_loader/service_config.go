package service_config_loader

import (
	"infer-microservices/pkg/config_loader/faiss_config"
	"infer-microservices/pkg/config_loader/model_config"
	"infer-microservices/pkg/config_loader/pipeline_config"
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
	serviceId         string                              `validate:"required,unique,min=4,max=10"` //dataid
	redisConfig       redis_config.RedisConfig            `validate:"required"`                     //redis conn info
	faissIndexConfigs faiss_config.FaissIndexConfigs      //index conn info
	modelsConfig      map[string]model_config.ModelConfig `validate:"required"` //model conn info
	pipelineCnfig     pipeline_config.PipelineConfig      `validate:"required"` //infer pipeline info
}

func init() {
}

// serviceId
func (s *ServiceConfig) setServiceId(serviceId string) {
	s.serviceId = serviceId
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
func (s *ServiceConfig) setModelsConfig(modelsConfig map[string]model_config.ModelConfig) {
	s.modelsConfig = modelsConfig
}

func (s *ServiceConfig) GetModelsConfig() map[string]model_config.ModelConfig {
	return s.modelsConfig
}

// piplineCnfig
func (s *ServiceConfig) setPipelineConfig(piplineCnfig pipeline_config.PipelineConfig) {
	s.pipelineCnfig = piplineCnfig
}

func (s *ServiceConfig) GetPipelineConfig() *pipeline_config.PipelineConfig {
	return &s.pipelineCnfig
}
