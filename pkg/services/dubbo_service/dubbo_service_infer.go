package server

import (
	"context"
	"errors"
	"fmt"
	config_loader "infer-microservices/pkg/config_loader"
	"infer-microservices/pkg/services/baseservice"

	"infer-microservices/pkg/logs"
	nacos "infer-microservices/pkg/nacos"
	"infer-microservices/pkg/services/io"
	"infer-microservices/pkg/utils"
	"time"

	_ "dubbo.apache.org/dubbo-go/v3/imports"
)

// extend from  baseservice
type DubboService struct {
	baseservice *baseservice.BaseService
}

func (s *DubboService) SetBaseService(baseservice *baseservice.BaseService) {
	s.baseservice = baseservice
}

func (s *DubboService) GetBaseService() *baseservice.BaseService {
	return s.baseservice
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

	for {
		select {
		case <-ctx.Done():
			switch ctx.Err() {
			case context.DeadlineExceeded:
				logs.Warn(requestId, time.Now(), "context timeout DeadlineExceeded.")
				return response, ctx.Err()
			case context.Canceled:
				logs.Warn(requestId, time.Now(), "context timeout Canceled.")
				return response, ctx.Err()
			}
		case responseCh := <-respCh:
			response = responseCh
			logs.Info(requestId, time.Now(), "response:", response)

			return response, nil
		}
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
	nacosConfig := nacosFactory.CreateNacosConfig(s.baseservice.GetNacosIp(), uint64(s.baseservice.GetNacosPort()), in)
	logs.Debug(requestId, time.Now(), "nacosConfig:", nacosConfig)

	nacosConfig.StartListenNacos()

	//infer
	ServiceConfig := config_loader.ServiceConfigs[in.GetDataId()]
	response_, err := s.baseservice.RecommenderInferHystrix(nil, "dubboServer", in, ServiceConfig)
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
