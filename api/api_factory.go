package api

import (
	dubbo_api "infer-microservices/api/dubbo_api/server"
	grpc_api "infer-microservices/api/grpc_api/server"
	rest_api "infer-microservices/api/rest_api/server"
	"infer-microservices/internal/flags"
	"infer-microservices/pkg/services"
)

var restPort uint
var grpcPort uint
var maxCpuNum int
var dubboConfFile string
var serviceFactory services.ServiceFactory

type ApiFactory struct {
}

type ServiceStartInterface interface {
	ServiceStart()
}

func init() {
	flagFactory := flags.FlagFactory{}
	//flagServiceConfig
	flagServiceConfig := flagFactory.CreateFlagServiceConfig()
	grpcPort = *flagServiceConfig.GetServiceRestPort()
	restPort = *flagServiceConfig.GetServiceGrpcPort()
	maxCpuNum = *flagServiceConfig.GetServiceMaxCpuNum()

	//flagDubbo
	flagDubbo := flagFactory.CreateFlagDubbo()
	dubboConfFile = *flagDubbo.GetDubboServiceFile()
}

// create dubbo server
func (f ApiFactory) CreateDubboServiceApi() *dubbo_api.DubboServiceApi {
	dubboServiceApi := new(dubbo_api.DubboServiceApi)
	dubboServiceApi.SetDubboConfFile(dubboConfFile)

	dubboService := serviceFactory.CreateDubboService()
	dubboServiceApi.SetDubboService(dubboService)

	return dubboServiceApi
}

// create grpc server
func (f ApiFactory) CreateGrpcServiceApi() *grpc_api.GrpcServiceApi {
	grpcServiceApi := new(grpc_api.GrpcServiceApi)
	grpcServiceApi.SetServicePort(grpcPort)
	grpcServiceApi.SetMaxCpuNum(maxCpuNum)

	grpcService := serviceFactory.CreateGrpcService()
	grpcServiceApi.SetGrpcService(grpcService)

	return grpcServiceApi
}

// create rest server
func (f ApiFactory) CreateRestServiceApi() *rest_api.HttpServiceApi {
	restServiceApi := new(rest_api.HttpServiceApi)
	restServiceApi.SetServicePort(restPort)
	restServiceApi.SetMaxCpuNum(maxCpuNum)

	restService := serviceFactory.CreateRestService()
	restServiceApi.SetRestService(restService)

	return restServiceApi
}
