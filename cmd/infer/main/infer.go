package main

import (
	"flag"
	"infer-microservices/api"
	"infer-microservices/internal/logs"
	"infer-microservices/pkg/infer_samples"
	"time"
)

//TODO:需要注意指针的滥用，协程里一个局部变量的改变是否会影响另一个协程中的变量。看看地址是否一样
//验证表明，局部变量地址不一样

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
		infer_samples.BloomClean(infer_samples.GetUserBloomFilterInstance())
		infer_samples.BloomClean(infer_samples.GetItemBloomFilterInstance())
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
	go infer_samples.WatchBloomConfig() //0 o'clock start service and load all users and all items into bloom filter.
	go resetBloom()                     //0 o'clock clean bloom filter, every 7 days.

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
