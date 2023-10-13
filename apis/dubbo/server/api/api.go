package api

import (
	"context"
	"infer-microservices/apis"

	"dubbo.apache.org/dubbo-go/v3/config"
	hessian "github.com/apache/dubbo-go-hessian2"
)

var (
	DubbogoInferServiceClient = &DubbogoInferService{}
)

func init() {
	//INFO: both input and output need to register.
	hessian.RegisterPOJO(&apis.RecRequest{})
	hessian.RegisterPOJO(&apis.RecResponse{})

	config.SetConsumerService(DubbogoInferServiceClient)
}

type DubbogoInferService struct {
	// define service func name.
	DubboRecommendServer func(ctx context.Context, req *apis.RecRequest) (*apis.RecResponse, error)
}

// refer : https://www.w3cschool.cn/dubbo/languages-golang-dubbo-go-30-quickstart-quickstart_dubbo.html
