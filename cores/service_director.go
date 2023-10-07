package cores

import (
	"infer-microservices/cores/service_config"
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

// init service conn from local file

type ServiceConfigDirector struct {
}

// TODO:一定要考虑多线程访问的更新问题
func (s *ServiceConfigDirector) ServiceConfigUpdateContainIndexDirector(domain string, dataId string,
	redisConfStr string, modelConfStr string, indexConfStr string) service_config.ServiceConfig {
	//load redis,faiss,model
	serviceConfigBuilder := service_config.ServiceConfigBuilder{} //{*service_config.GetServiceConfigInstance()}, 私有方法，改成set
	builder := serviceConfigBuilder.RedisClientBuilder(domain, dataId, redisConfStr).FaissClientBuilder(indexConfStr).ModelClientBuilder(modelConfStr)

	return builder.GetServiceConfig()
}

// TODO:此处无返回
func (s *ServiceConfigDirector) ServiceConfigUpdaterNotContainIndexDirector(domain string, dataId string,
	redisConfStr string, modelConfStr string) service_config.ServiceConfig {

	//load redis,model
	serviceConfigBuilder := service_config.ServiceConfigBuilder{} //{*service_config.GetServiceConfigInstance()}, 私有方法，改成set
	builder := serviceConfigBuilder.RedisClientBuilder(domain, dataId, redisConfStr).ModelClientBuilder(modelConfStr)

	return builder.GetServiceConfig()
}
