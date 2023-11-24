package main

import (
	"flag"
	"infer-microservices/api"
	"infer-microservices/internal"
	"infer-microservices/pkg/logs"
	"infer-microservices/pkg/model/basemodel"
	"time"
)

var apiFactory api.ApiFactory

func init() {
	apiFactory = api.ApiFactory{}
}

// reset bloomfilter
func resetBloom() {
	//TODO: Cuckoo Filter better
	var ticker *time.Ticker = time.NewTicker(7 * 24 * time.Hour) //every 7 days clean bloom filter
	defer ticker.Stop()

	for t := range ticker.C {
		logs.Info("Start to reset bloom filter.", t)
		internal.CleanBloom(internal.GetUserBloomFilterInstance())
		internal.CleanBloom(internal.GetItemBloomFilterInstance())
	}
}

// start  service
func serviceStart(serviceApi api.ServiceStartInterface) {
	logs.Info("starting dubbo servivce...")
	serviceApi.ServiceStart()
	logs.Info("successed start dubbo servivce.")
}

func main() {
	//init
	flag.Parse()
	logs.InitLog()

	//watch and reset bloom fliter
	go basemodel.WatchBloomConfig() //0 o'clock start service and load all users and all items into bloom filter.
	go resetBloom()                 //0 o'clock clean bloom filter, every 7 days.

	//start services.
	dubboServiceApi := apiFactory.CreateDubboServiceApi()
	grpcServiceApi := apiFactory.CreateGrpcServiceApi()
	//restServiceHttpApi := apiFactory.CreateRestServiceHttpApi()  // @deprecated
	restServiceEchoApi := apiFactory.CreateRestServiceEchoApi()

	go serviceStart(dubboServiceApi)
	go serviceStart(grpcServiceApi)
	//go serviceStart(restServiceApi)  // @deprecated
	go serviceStart(restServiceEchoApi)

	select {}
}
