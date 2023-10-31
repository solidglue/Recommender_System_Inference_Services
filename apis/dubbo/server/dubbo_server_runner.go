package server

import (
	"infer-microservices/apis/io"

	"dubbo.apache.org/dubbo-go/v3/config"
	hessian "github.com/apache/dubbo-go-hessian2"
)

func init() {
	//regisger dubbo service.
	hessian.RegisterPOJO(&io.RecRequest{})
	hessian.RegisterPOJO(&io.RecResponse{})
	config.SetProviderService(&DubbogoInferService{})
}

func DubboServerRunner(ipAddr string, port uint64, dubboConfFile string) {
	ipAddr_ = ipAddr
	port_ = port
	if err := config.Load(config.WithPath(dubboConfFile)); err != nil {
		panic(err)
	}
}
