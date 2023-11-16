package services

import (
	"infer-microservices/internal/flags"
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

// create dubbo server
func (f ServiceFactory) CreateDubboService() *dubbo_service.DubboService {
	dubboService := new(dubbo_service.DubboService)
	dubboService.SetNacosIp(nacosIp)
	dubboService.SetNacosPort(uint(nacosPort))
	dubboService.SetLowerRankNum(lowerRankNum)
	dubboService.SetLowerRecallNum(lowerRecallNum)

	return dubboService
}

// create grpc server
func (f ServiceFactory) CreateGrpcService() *grpc_service.GrpcService {
	grpcService := new(grpc_service.GrpcService)
	grpcService.SetNacosIp(nacosIp)
	grpcService.SetNacosPort(uint(nacosPort))
	grpcService.SetSkywalkingWeatherOpen(false)
	grpcService.SetLowerRankNum(lowerRankNum)
	grpcService.SetLowerRecallNum(lowerRecallNum)

	return grpcService
}

// create rest server
func (f ServiceFactory) CreateRestService() *rest_service.HttpService {
	restService := new(rest_service.HttpService)
	restService.SetNacosIp(nacosIp)
	restService.SetNacosPort(uint(nacosPort))
	restService.SetSkywalkingWeatherOpen(skywalkingWeatherOpen)
	restService.SetSkywalkingIp(skywalkingIp)
	restService.SetSkywalkingPort(uint(skywalkingPort))
	restService.SetSkywalkingServerName(skywalkingServerName)
	restService.SetLowerRankNum(lowerRankNum)
	restService.SetLowerRecallNum(lowerRecallNum)

	return restService
}
