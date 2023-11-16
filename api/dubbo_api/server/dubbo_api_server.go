package server

import (
	"context"
	dubbo_service "infer-microservices/pkg/services/dubbo_service"
	"infer-microservices/pkg/services/io"

	"dubbo.apache.org/dubbo-go/v3/config"
	hessian "github.com/apache/dubbo-go-hessian2"
)

var (
	DubboServiceApiClient = &DubboServiceApi{}
)

type DubboServiceApi struct {
	dubboConfFile string
	dubboService  *dubbo_service.DubboService
	// define service func name.
	RecommenderInfer func(ctx context.Context, req *io.RecRequest) (*io.RecResponse, error)
}

func init() {
	//regisger dubbo service.
	//INFO: both input and output need to register.
	hessian.RegisterPOJO(&io.RecRequest{})
	hessian.RegisterPOJO(&io.RecResponse{})

	config.SetConsumerService(DubboServiceApiClient)

}

func (s *DubboServiceApi) SetDubboConfFile(dubboConfFile string) {
	s.dubboConfFile = dubboConfFile
}

func (s *DubboServiceApi) SetDubboService(dubboService *dubbo_service.DubboService) {
	s.dubboService = dubboService
}

// @implement start infertace
func (s *DubboServiceApi) ServiceStart() {
	config.SetProviderService(s)
	if err := config.Load(config.WithPath(s.dubboConfFile)); err != nil {
		panic(err)
	}
}
