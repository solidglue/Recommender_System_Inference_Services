package main

import (
	"flag"
	"infer-microservices/apis"
	"infer-microservices/common"
	"infer-microservices/utils/logs"
	"time"
)

var apiFactory apis.ApiFactory
var apiServer apis.ServerStartInterface //INFO:apis implement ServerStartInterface

func init() {
	apiFactory = apis.ApiFactory{}
}

// reset bloomfilter
func resetBloom() {
	var ticker *time.Ticker = time.NewTicker(7 * 24 * time.Hour) //every 7 days clean bloom filter
	defer ticker.Stop()

	for t := range ticker.C {
		logs.Info("Start to reset bloom filter.", t)
		common.CleanBloom(common.GetUserBloomFilterInstance())
		common.CleanBloom(common.GetItemBloomFilterInstance())
	}
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

	go common.WatchBloomConfig() //0 o'clock start service and load all users and all items into bloom filter.

	go dubboServiceStart()
	go grpcServiceStart()
	go restfulServiceStart()

	go resetBloom() //0 o'clock clean bloom filter, every 7 days.

	select {}
}
