package server

import (
	"errors"
	"fmt"
	"infer-microservices/cores/model"
	"infer-microservices/cores/nacos_config_listener"
	"infer-microservices/cores/service_config_loader"
	"strings"
	"sync"
	"time"

	grpc_api "infer-microservices/apis/grpc/server/api_gogofaster"
	"infer-microservices/apis/io"

	"infer-microservices/utils/logs"

	"github.com/afex/hystrix-go/hystrix"
	"golang.org/x/net/context"
)

var inferWg sync.WaitGroup

// INFO: implement grpc func which defined by proto.
func (g *GrpcServer) GrpcRecommendServer(ctx context.Context, in *grpc_api.RecommendRequest) (*grpc_api.RecommendResponse, error) {
	//INFO: set timeout by context, degraded service by hystix.
	response := &grpc_api.RecommendResponse{
		Code: 404,
	}

	ctx, cancelFunc := context.WithTimeout(ctx, time.Millisecond*150)
	defer cancelFunc()

	respCh := make(chan *grpc_api.RecommendResponse, 100)
	go g.grpcRecommenderServerContext(ctx, in, respCh)
	for {
		select {
		case <-ctx.Done():
			switch ctx.Err() {
			case context.DeadlineExceeded:
				return response, ctx.Err()
			case context.Canceled:
				return response, ctx.Err()
			}
		case responseCh := <-respCh:
			response = responseCh
			return response, nil
		}
	}
}

func (s *GrpcServer) grpcRecommenderServerContext(ctx context.Context, in *grpc_api.RecommendRequest, respCh chan *grpc_api.RecommendResponse) {
	defer func() {
		if info := recover(); info != nil {
			fmt.Println("panic", info)
		} //else {
		//fmt.Println("")
		//}
	}()

	response := &grpc_api.RecommendResponse{
		Code: 404,
	}

	request := convertGrpcRequestToRecRequest(in)
	//check input
	checkStatus := request.Check()
	if !checkStatus {
		err := errors.New("input check failed")
		logs.Error(err)
		panic(err)
	}

	//nacos listen
	nacosFactory := nacos_config_listener.NacosFactory{}
	nacosConfig := nacosFactory.CreateNacosConfig(s.nacosIp, uint64(s.nacosPort), &request)
	nacosConfig.StartListenNacos()

	//infer
	ServiceConfig := service_config_loader.ServiceConfigs[in.GetDataId()]
	response_, err := s.grpcHystrixServer("grpcServer", &request, ServiceConfig)
	if err != nil {
		response.Message = fmt.Sprintf("%s", err)
		panic(err)
	} else {
		response = response_
	}
	respCh <- response
}

func convertGrpcRequestToRecRequest(in *grpc_api.RecommendRequest) io.RecRequest {
	request := io.RecRequest{}
	request.SetDataId(in.GetDataId())
	request.SetGroupId(in.GetGroupId())
	request.SetNamespaceId(in.GetNamespace())
	request.SetUserId(in.UserId)
	request.SetRecallNum(in.RecallNum)
	request.SetItemList(in.ItemList.Value)

	return request
}

func (s *GrpcServer) grpcHystrixServer(serverName string, in *io.RecRequest, ServiceConfig *service_config_loader.ServiceConfig) (*grpc_api.RecommendResponse, error) {
	response := &grpc_api.RecommendResponse{
		Code: 404,
	}

	hystrixErr := hystrix.Do(serverName, func() error {
		// request recall / rank func.
		response_, err := s.recommenderInfer(in, ServiceConfig)
		if err != nil {
			logs.Error(err)
			return err
		} else {
			response = response_
		}

		return nil
	}, func(err error) error {
		//INFO: do this when services are timeout (hystrix timeout).
		// less items and simple model.
		itemList := in.GetItemList()
		in.SetRecallNum(int32(s.lowerRecallNum))
		in.SetItemList(itemList[:s.lowerRankNum])
		response_, err_ := s.recommenderInferReduce(in, ServiceConfig)

		if err_ != nil {
			logs.Error(err_)
			return err_
		} else {
			response = response_
		}

		return nil
	})

	if hystrixErr != nil {
		return response, hystrixErr
	}

	return response, nil
}

func (s *GrpcServer) recommenderInfer(in *io.RecRequest, ServiceConfig *service_config_loader.ServiceConfig) (*grpc_api.RecommendResponse, error) {
	response := &grpc_api.RecommendResponse{
		Code: 404,
	}

	//build model by model_factory
	modelName := in.GetModelType()
	if modelName != "" {
		modelName = strings.ToLower(modelName)
	}

	//strategy pattern
	modelfactory := model.ModelStrategyFactory{}
	modelStrategy := modelfactory.CreateModelStrategy(modelName, in, ServiceConfig)
	modelStrategyContext := model.ModelStrategyContext{}
	modelStrategyContext.SetModelStrategy(modelStrategy)

	var err error
	var responseTf map[string]interface{}

	if s.skywalkingWeatherOpen {
		responseTf, err = modelStrategyContext.ModelInferSkywalking(nil)
	} else {
		responseTf, err = modelStrategyContext.ModelInferNoSkywalking(nil)
	}
	if err != nil {
		logs.Error(err)
		return response, err
	}

	result := make([]*grpc_api.ItemInfo, 0)
	resultList := responseTf["data"].([]map[string]interface{})
	rankCh := make(chan *grpc_api.ItemInfo, len(resultList))
	for i := 0; i < len(resultList); i++ {
		inferWg.Add(1)
		go formatGrpcResponse(resultList[i], rankCh)
	}

	inferWg.Wait()
	close(rankCh)
	for itemScore := range rankCh {
		result = append(result, itemScore)
	}

	response = &grpc_api.RecommendResponse{
		Code:    200,
		Message: "success",
		Data: &grpc_api.ItemInfoList{
			Iteminfo_: result,
		},
	}

	return response, nil
}

func (s *GrpcServer) recommenderInferReduce(in *io.RecRequest, ServiceConfig *service_config_loader.ServiceConfig) (*grpc_api.RecommendResponse, error) {
	response := &grpc_api.RecommendResponse{
		Code: 404,
	}

	//build model by model_factory
	// modelName := in.GetModelType()
	// if modelName != "" {
	// 	modelName = strings.ToLower(modelName)
	// }

	modelName := "fm"

	//strategy pattern
	modelfactory := model.ModelStrategyFactory{}
	modelStrategy := modelfactory.CreateModelStrategy(modelName, in, ServiceConfig)
	modelStrategyContext := model.ModelStrategyContext{}
	modelStrategyContext.SetModelStrategy(modelStrategy)

	responseTf, err := modelStrategyContext.ModelInferSkywalking(nil)
	if err != nil {
		logs.Error(err)
		return response, err
	}

	result := make([]*grpc_api.ItemInfo, 0)
	resultList := responseTf["data"].([]map[string]interface{})
	rankCh := make(chan *grpc_api.ItemInfo, len(resultList))
	for i := 0; i < len(resultList); i++ {
		inferWg.Add(1)
		go formatGrpcResponse(resultList[i], rankCh)
	}

	inferWg.Wait()
	close(rankCh)
	for itemScore := range rankCh {
		result = append(result, itemScore)
	}

	response = &grpc_api.RecommendResponse{
		Code:    200,
		Message: "success",
		Data: &grpc_api.ItemInfoList{
			Iteminfo_: result,
		},
	}

	return response, nil
}

func formatGrpcResponse(itemScore map[string]interface{}, rankCh chan *grpc_api.ItemInfo) {
	defer inferWg.Done()

	itemid := itemScore["itemid"].(string)
	score := float32(itemScore["score"].(float64))
	itemInfo := &grpc_api.ItemInfo{
		Itemid: itemid,
		Score:  score,
	}

	rankCh <- itemInfo
}
