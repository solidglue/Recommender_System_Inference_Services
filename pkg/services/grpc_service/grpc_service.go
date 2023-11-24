package grpc_service

import (
	"errors"
	"fmt"
	config_loader "infer-microservices/pkg/config_loader"
	"infer-microservices/pkg/utils"
	"time"

	"infer-microservices/pkg/logs"
	nacos "infer-microservices/pkg/nacos"
	"infer-microservices/pkg/services/baseservice"
	"infer-microservices/pkg/services/io"

	"golang.org/x/net/context"
)

// extend from  baseservice
type GrpcService struct {
	baseservice *baseservice.BaseService
}

func (s *GrpcService) SetBaseService(baseservice *baseservice.BaseService) {
	s.baseservice = baseservice
}

func (s *GrpcService) GetBaseService() *baseservice.BaseService {
	return s.baseservice
}

// INFO: implement grpc func which defined by proto.
func (s *GrpcService) RecommenderInfer(ctx context.Context, in *RecommendRequest) (*RecommendResponse, error) {
	//INFO: set timeout by context, degraded service by hystix.
	response := &RecommendResponse{
		Code: 404,
	}
	request := convertGrpcRequestToRecRequest(in)
	requestId := utils.CreateRequestId(&request)
	logs.Debug(requestId, time.Now(), "RecRequest:", requestId)

	ctx, cancelFunc := context.WithTimeout(ctx, time.Millisecond*100)
	defer cancelFunc()

	respCh := make(chan *RecommendResponse, 100)
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

func (s *GrpcService) recommenderInferContext(ctx context.Context, in *RecommendRequest, respCh chan *RecommendResponse) {
	defer func() {
		if info := recover(); info != nil {
			logs.Fatal("panic", info)
		} //else {
		//fmt.Println("")
		//}
	}()

	response := &RecommendResponse{
		Code: 404,
	}
	request := convertGrpcRequestToRecRequest(in)
	requestId := utils.CreateRequestId(&request)

	//check input
	checkStatus := request.Check()
	if !checkStatus {
		err := errors.New("input check failed")
		logs.Error(requestId, time.Now(), err)
		panic(err)
	}

	//nacos listen
	nacosFactory := nacos.NacosFactory{}
	nacosConfig := nacosFactory.CreateNacosConfig(s.baseservice.GetNacosIp(), uint64(s.baseservice.GetNacosPort()), &request)
	logs.Debug(requestId, time.Now(), "nacosConfig:", nacosConfig)

	nacosConfig.StartListenNacos()

	//infer
	ServiceConfig := config_loader.GetServiceConfigs()[in.GetDataId()]
	response_, err := s.baseservice.RecommenderInferHystrix(nil, "GrpcService", &request, ServiceConfig)
	if err != nil {
		response.Message = fmt.Sprintf("%s", err)
		panic(err)
	} else {
		response = &RecommendResponse{
			Code:    response_["code"].(int32),
			Message: response_["message"].(string),
			Data: &ItemInfoList{
				Iteminfo_: response_["data"].([]*ItemInfo),
			},
		}

	}
	respCh <- response
}

func convertGrpcRequestToRecRequest(in *RecommendRequest) io.RecRequest {
	request := io.RecRequest{}
	request.SetDataId(in.GetDataId())
	request.SetGroupId(in.GetGroupId())
	request.SetNamespaceId(in.GetNamespace())
	request.SetUserId(in.UserId)
	request.SetRecallNum(in.RecallNum)
	request.SetItemList(in.ItemList.Value)

	return request
}
