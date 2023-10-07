package main

import (
	"flag"
	dubbo_server "infer-microservices/apis/dubbo/server"
	grpc_server "infer-microservices/apis/grpc/server"
	rest_server "infer-microservices/apis/rest/server"
	"infer-microservices/common/flags"
	"infer-microservices/utils/logs"
)

var dubboConf string
var nacosIp string
var nacosPort uint64

func init() {
	flagFactory := flags.FlagFactory{}
	flagDubbo := flagFactory.FlagDubboFactory()
	flagNacos := flagFactory.FlagNacosFactory()

	dubboConf = *flagDubbo.GetDubboServiceFile()
	nacosIp = *flagNacos.GetNacosIp()
	nacosPort = uint64(*flagNacos.GetNacosPort())
}

func restfulServiceStart() {
	rest_server.NacosIP = nacosIp
	rest_server.NacosPort = nacosPort
	logs.Info("starting rest servivce...")
	rest_server.RestServerRunner()
	logs.Info("finished start servivce.")
}

func grpcServiceStart() {
	logs.Info("starting grpc servivce...")
	grpc_server.GrpcServerRunner(nacosIp, nacosPort)
	logs.Info("finished start servivce.")
}

func dubboServiceStart() {
	logs.Info("starting dubbo servivce...")
	dubbo_server.DubboServerRunner(nacosIp, nacosPort, dubboConf)
	logs.Info("finished dubbo servivce.")
}

func main() {
	flag.Parse()
	logs.InitLog()

	go dubboServiceStart()
	go grpcServiceStart()
	go restfulServiceStart()

	select {}
}
