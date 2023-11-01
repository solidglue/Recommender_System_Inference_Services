package main

import (
	"flag"
	"infer-microservices/apis"
	"infer-microservices/utils/logs"
)

var apiFactory apis.ApiFactory
var apiServer apis.ServerStartInterface //INFO:apis implement ServerStartInterface

func init() {
	apiFactory = apis.ApiFactory{}
}

// start dubbo service
func dubboServiceStart() {
	apiServer = apiFactory.CreateDubboServer()
	logs.Info("starting dubbo servivce...")
	apiServer.ServerStart()
	logs.Info("successed start dubbo servivce.")
}

// start grpc service
func grpcServiceStart() {
	apiServer = apiFactory.CreateGrpcServer()
	logs.Info("starting grpc servivce...")
	apiServer.ServerStart()
	logs.Info("successed start grpc servivce.")
}

// start rest service
func restfulServiceStart() {
	apiServer = apiFactory.CreateRestServer()
	logs.Info("starting rest servivce...")
	apiServer.ServerStart()
	logs.Info("successed start rest servivce.")
}

func main() {
	flag.Parse()
	logs.InitLog()

	go dubboServiceStart()
	go grpcServiceStart()
	go restfulServiceStart()

	select {}
}
