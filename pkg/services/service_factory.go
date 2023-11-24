package services

import (
	"infer-microservices/internal/flags"
	"infer-microservices/pkg/services/baseservice"
	dubbo_service "infer-microservices/pkg/services/dubbo_service"
	grpc_service "infer-microservices/pkg/services/grpc_service"
	rest_service "infer-microservices/pkg/services/rest_service"
)

var skywalkingWeatherOpen bool
var skywalkingIp string
var skywalkingPort int
var skywalkingServerName string
var nacosIp string
var nacosPort uint64
var lowerRankNum int
var lowerRecallNum int

type RecommenderInferInterface interface {
	RecommenderInfer()
}

type ServiceFactory struct {
}

func init() {
	flagFactory := flags.FlagFactory{}
	//flagSkywalking
	flagSkywalking := flagFactory.CreateFlagSkywalking()
	skywalkingWeatherOpen = *flagSkywalking.GetSkywalkingWhetheropen()
	skywalkingIp = *flagSkywalking.GetSkywalkingIp()
	skywalkingPort = *flagSkywalking.GetSkywalkingPort()
	skywalkingServerName = *flagSkywalking.GetSkywalkingServername()

	//flagHystrix
	flagHystrix := flagFactory.CreateFlagHystrix()
	lowerRecallNum = *flagHystrix.GetHystrixLowerRecallNum()
	lowerRankNum = *flagHystrix.GetHystrixLowerRankNum()
}

// create base server
func (f ServiceFactory) createBaseServiceSkywalking() *baseservice.BaseService {
	baseService := new(baseservice.BaseService)
	baseService.SetNacosIp(nacosIp)
	baseService.SetNacosPort(uint(nacosPort))
	baseService.SetSkywalkingWeatherOpen(skywalkingWeatherOpen)
	baseService.SetSkywalkingIp(skywalkingIp)
	baseService.SetSkywalkingPort(uint(skywalkingPort))
	baseService.SetSkywalkingServerName(skywalkingServerName)
	baseService.SetLowerRankNum(lowerRankNum)
	baseService.SetLowerRecallNum(lowerRecallNum)

	return baseService
}

// create base server
func (f ServiceFactory) createBaseServiceNoSkywalking() *baseservice.BaseService {
	baseService := new(baseservice.BaseService)
	baseService.SetNacosIp(nacosIp)
	baseService.SetNacosPort(uint(nacosPort))
	baseService.SetSkywalkingWeatherOpen(false)
	baseService.SetLowerRankNum(lowerRankNum)
	baseService.SetLowerRecallNum(lowerRecallNum)

	return baseService
}

// create dubbo server
func (f ServiceFactory) CreateDubboService() *dubbo_service.DubboService {
	dubboService := new(dubbo_service.DubboService)
	dubboService.SetBaseService(f.createBaseServiceNoSkywalking())

	return dubboService
}

// create grpc server
func (f ServiceFactory) CreateGrpcService() *grpc_service.GrpcService {
	grpcService := new(grpc_service.GrpcService)
	grpcService.SetBaseService(f.createBaseServiceNoSkywalking())

	return grpcService
}

// create http server
// @deprecated
func (f ServiceFactory) CreateHttpService() *rest_service.HttpService {
	httpService := new(rest_service.HttpService)
	httpService.SetBaseService(f.createBaseServiceSkywalking())

	return httpService
}

// create echo server
func (f ServiceFactory) CreateEchoService() *rest_service.EchoService {
	echoService := new(rest_service.EchoService)
	echoService.SetBaseService(f.createBaseServiceSkywalking())

	return echoService
}
