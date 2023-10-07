package service_config

import (
	"infer-microservices/common/flags"
	"infer-microservices/cores/faiss"
	"infer-microservices/cores/model"
	"infer-microservices/cores/redis_config"
)

// TODO:改成单例模式
// TODO:viper配置解决方案，可以自动加载监听文件，无需重启服务
// TODO:validator，验证传入的json
// TODO:日志系统替换为zap（高性能log）
// TODO:使用Asynq异步写日志，比较影响主流程性能
// TODO:用casbin进行权限管控，测试、开发、生产等

//TODO:都定义loader或infer接口，类去实现，并用接口调用。参数是具体的类。结合工厂模式。
//小配置用工厂模式？大配置用建造者模式，连续build？
//补充封装get set

//todo：配置文件大类用建造者模式，主要用于创建复杂对象 + 单例模式（单例模式放在产品类里？）。其它类用工厂模式，主要用于统一管控new对象

//TODO:加载model、redis、faiss用工厂模式
//TODO:生成模型配置用建造模型（faiss非必须）
//TODO:服务配置是单例模式

// TODO:加锁
var serviceConfFile string
var serviceConfigInstance *ServiceConfig

//TODO:ServiceConfig与dataid相关，无需单例。load的时候把dataid加上，取实例开放给dubboserver方法，或工厂模式

type ServiceConfig struct {
	serviceId        string                   //dataid
	redisClient      redis_config.RedisClient //redis conn info
	faissIndexClient faiss.FaissIndexClient   //index conn info
	modelClient      model.ModelClient        //model conn info
}

func init() {
	//serviceConfFile = *flags.Service_start_file

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
