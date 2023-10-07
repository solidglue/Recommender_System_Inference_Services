package main

import (
	"flag"
	"fmt"
	dubbo_server "infer-microservices/apis/dubbo/server"
	grpc_server "infer-microservices/apis/grpc/server"
	rest_server "infer-microservices/apis/rest/server"
	"infer-microservices/common/flags"
	"infer-microservices/utils/logs"
)

var dubboConf string
var nacosIp string
var nacosPort uint64

//TODO:lock和指针梳理。20230804

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
	fmt.Println("starting rest servivce...")
	rest_server.RestServerRunner()
	fmt.Println("finished start servivce.")
}

func grpcServiceStart() {

	//启动grpc
	fmt.Println("starting grpc servivce...")
	grpc_server.GrpcServerRunner(nacosIp, nacosPort) //此处不能用go协程，否则立刻退出/完成执行
	fmt.Println("finished start servivce.")

}

func dubboServiceStart() {

	fmt.Println("starting dubbo servivce...")
	dubbo_server.DubboServerRunner(nacosIp, nacosPort, dubboConf)
	fmt.Println("finished dubbo servivce.")

}

func main() {

	flag.Parse()
	logs.InitLog()

	go dubboServiceStart()
	go grpcServiceStart()
	go restfulServiceStart()

	select {}

}
