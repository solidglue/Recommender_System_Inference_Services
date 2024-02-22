package dubbo_service

import (
	"context"
	"errors"
	"fmt"
	config_loader "infer-microservices/pkg/config_loader"
	"infer-microservices/pkg/infer_services/base_service"

	"infer-microservices/internal/logs"
	"infer-microservices/internal/utils"
	"infer-microservices/pkg/infer_services/io"
	nacos "infer-microservices/pkg/nacos"
	"time"

	_ "dubbo.apache.org/dubbo-go/v3/imports"
)

// extend from  baseservice
type DubboService struct {
	baseService *base_service.BaseService
}

func (s *DubboService) SetBaseService(baseService *base_service.BaseService) {
	s.baseService = baseService
}

func (s *DubboService) GetBaseService() *base_service.BaseService {
	return s.baseService
}

//INFO:DONT REMOVE.  JAVA request service need it.
// // MethodMapper mapper upper func name to lower func name ,for java request.
// func (s *InferDubbogoService) MethodMapper() map[string]string {
// 	return map[string]string{
// 		"DubboRecommendServer": "dubboRecommendServer",
// 	}
// }

// Implement interface methods.
func (s *DubboService) RecommenderInfer(ctx context.Context, in *io.RecRequest) (*io.RecResponse, error) {
	response := &io.RecResponse{}
	response.SetCode(404)
	requestId := utils.CreateRequestId(in)
	logs.Debug(requestId, time.Now(), "RecRequest:", in)

	//check input
	checkStatus := in.Check()
	if !checkStatus {
		err := errors.New("input check failed")
		logs.Error(requestId, time.Now(), err)
		return response, err
	}

	//INFO: set timeout by context, degraded service by hystix.
	ctx, cancelFunc := context.WithTimeout(ctx, time.Millisecond*100)
	defer cancelFunc()

	respCh := make(chan *io.RecResponse, 100)
	go s.recommenderInferContext(ctx, in, respCh)

	select {
	case <-ctx.Done():
		switch ctx.Err() {
		case context.DeadlineExceeded:
			logs.Error(requestId, time.Now(), ctx.Err())
		case context.Canceled:
			logs.Error(requestId, time.Now(), ctx.Err())
		}
		return response, ctx.Err()
	case responseCh := <-respCh:
		response = responseCh
		logs.Info(requestId, time.Now(), "response:", response)
		return response, nil
	}
}

func (s *DubboService) recommenderInferContext(ctx context.Context, in *io.RecRequest, respCh chan *io.RecResponse) {
	defer func() {
		if info := recover(); info != nil {
			logs.Fatal("panic", info)
		} //else {
		//  logs.Info("finish.")
		//}
	}()

	response := &io.RecResponse{}
	response.SetCode(404)
	requestId := utils.CreateRequestId(in)

	//nacos listen
	nacosFactory := nacos.NacosFactory{}
	nacosConfig := nacosFactory.CreateNacosConfig(s.baseService.GetNacosIp(), uint64(s.baseService.GetNacosPort()), in)
	logs.Debug(requestId, time.Now(), "nacosConfig:", nacosConfig)

	nacosConfig.StartListenNacos()

	//infer
	ServiceConfig := config_loader.GetServiceConfigs()[in.GetDataId()]
	response_, err := s.baseService.RecommenderInferHystrix(nil, "dubboServer", in, ServiceConfig)
	if err != nil || len(response_) == 0 {
		response.SetMessage(fmt.Sprintf("%s", err))
		panic(err)
	} else {
		response.SetCode(response_["code"].(int))
		response.SetMessage(response_["message"].(string))
		response.SetData(response_["data"].([]string))
	}

	respCh <- response
}
