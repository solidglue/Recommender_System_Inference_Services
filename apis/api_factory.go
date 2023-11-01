package apis

import (
	dubbo_server "infer-microservices/apis/dubbo/server"
	grpc_server "infer-microservices/apis/grpc/server"
	rest_server "infer-microservices/apis/rest/server"
	"infer-microservices/common/flags"
)

var serverPort uint
var maxCpuNum int
var skywalkingWeatherOpen bool
var skywalkingIp string
var skywalkingPort int
var skywalkingServerName string
var nacosIp string
var nacosPort uint64
var lowerRankNum int
var lowerRecallNum int
var dubboConfFile string

type ApiFactory struct {
}

type ServerStartInterface interface {
	ServerStart()
}

func init() {
	flagFactory := flags.FlagFactory{}
	//flagServiceConfig
	flagServiceConfig := flagFactory.CreateFlagServiceConfig()
	serverPort = *flagServiceConfig.GetServiceRestPort()
	maxCpuNum = *flagServiceConfig.GetServiceMaxCpuNum()

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

	//flagDubbo
	flagDubbo := flagFactory.CreateFlagDubbo()
	dubboConfFile = *flagDubbo.GetDubboServiceFile()
}

// create dubbo server
func (f ApiFactory) CreateDubboServer() *dubbo_server.DubboServer {
	dubboServer := new(dubbo_server.DubboServer)
	dubboServer.SetDubboConfFile(dubboConfFile)
	dubboServer.SetNacosIp(nacosIp)
	dubboServer.SetNacosPort(uint(nacosPort))
	dubboServer.SetLowerRankNum(lowerRankNum)
	dubboServer.SetLowerRecallNum(lowerRecallNum)

	return dubboServer
}

// create grpc server
func (f ApiFactory) CreateGrpcServer() *grpc_server.GrpcServer {
	grpcServer := new(grpc_server.GrpcServer)
	grpcServer.SetServerPort(serverPort)
	grpcServer.SetNacosIp(nacosIp)
	grpcServer.SetNacosPort(uint(nacosPort))
	grpcServer.SetMaxCpuNum(maxCpuNum)
	grpcServer.SetSkywalkingWeatherOpen(false)
	grpcServer.SetLowerRankNum(lowerRankNum)
	grpcServer.SetLowerRecallNum(lowerRecallNum)

	return grpcServer
}

// create rest server
func (f ApiFactory) CreateRestServer() *rest_server.HttpServer {
	restServer := new(rest_server.HttpServer)
	restServer.SetServerPort(serverPort)
	restServer.SetNacosIp(nacosIp)
	restServer.SetNacosPort(uint(nacosPort))
	restServer.SetMaxCpuNum(maxCpuNum)
	restServer.SetSkywalkingWeatherOpen(skywalkingWeatherOpen)
	restServer.SetSkywalkingIp(skywalkingIp)
	restServer.SetSkywalkingPort(uint(skywalkingPort))
	restServer.SetSkywalkingServerName(skywalkingServerName)
	restServer.SetLowerRankNum(lowerRankNum)
	restServer.SetLowerRecallNum(lowerRecallNum)

	return restServer
}
