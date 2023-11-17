package grpc_service

import (
	"errors"
	"fmt"
	config_loader "infer-microservices/pkg/config_loader"
	"infer-microservices/pkg/model"
	"infer-microservices/pkg/utils"
	"strings"
	"sync"
	"time"

	"infer-microservices/pkg/logs"
	nacos "infer-microservices/pkg/nacos"
	"infer-microservices/pkg/services/io"

	"github.com/afex/hystrix-go/hystrix"
	"golang.org/x/net/context"
)

var inferWg sync.WaitGroup

type GrpcService struct {
	nacosIp               string
	nacosPort             uint
	skywalkingWeatherOpen bool
	lowerRankNum          int
	lowerRecallNum        int
}

// set func
func (s *GrpcService) SetNacosIp(nacosIp string) {
	s.nacosIp = nacosIp
}

func (s *GrpcService) SetNacosPort(nacosPort uint) {
	s.nacosPort = nacosPort
}

func (s *GrpcService) SetSkywalkingWeatherOpen(skywalkingWeatherOpen bool) {
	s.skywalkingWeatherOpen = skywalkingWeatherOpen
}

func (s *GrpcService) SetLowerRankNum(lowerRankNum int) {
	s.lowerRankNum = lowerRankNum
}

func (s *GrpcService) SetLowerRecallNum(lowerRecallNum int) {
	s.lowerRecallNum = lowerRecallNum
}

// INFO: implement grpc func which defined by proto.
func (g *GrpcService) RecommenderInfer(ctx context.Context, in *RecommendRequest) (*RecommendResponse, error) {
	//INFO: set timeout by context, degraded service by hystix.
	response := &RecommendResponse{
		Code: 404,
	}
	request := convertGrpcRequestToRecRequest(in)
	requestId := utils.CreateRequestId(&request)
	logs.Debug(requestId, time.Now(), "RecRequest:", requestId)

	ctx, cancelFunc := context.WithTimeout(ctx, time.Millisecond*150)
	defer cancelFunc()

	respCh := make(chan *RecommendResponse, 100)
	go g.recommenderInferContext(ctx, in, respCh)
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
	nacosConfig := nacosFactory.CreateNacosConfig(s.nacosIp, uint64(s.nacosPort), &request)
	logs.Debug(requestId, time.Now(), "nacosConfig:", nacosConfig)

	nacosConfig.StartListenNacos()

	//infer
	ServiceConfig := config_loader.ServiceConfigs[in.GetDataId()]
	response_, err := s.recommenderInferHystrix("GrpcService", &request, ServiceConfig)
	if err != nil {
		response.Message = fmt.Sprintf("%s", err)
		panic(err)
	} else {
		response = response_
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

func (s *GrpcService) recommenderInferHystrix(serverName string, in *io.RecRequest, ServiceConfig *config_loader.ServiceConfig) (*RecommendResponse, error) {
	response := &RecommendResponse{
		Code: 404,
	}
	requestId := utils.CreateRequestId(in)

	hystrixErr := hystrix.Do(serverName, func() error {
		// request recall / rank func.
		response_, err := s.modelInfer(in, ServiceConfig)
		if err != nil {
			logs.Error(requestId, time.Now(), err)
			return err
		} else {
			response = response_
		}

		return err
	}, func(err error) error {
		//INFO: do this when services are timeout (hystrix timeout).
		// less items and simple model.
		itemList := in.GetItemList()
		in.SetRecallNum(int32(s.lowerRecallNum))
		in.SetItemList(itemList[:s.lowerRankNum])
		response_, err_ := s.modelInferReduce(in, ServiceConfig)

		if err_ != nil {
			logs.Error(requestId, time.Now(), err_)
			return err_
		} else {
			response = response_
		}

		return err
	})

	if hystrixErr != nil {
		return response, hystrixErr
	}

	return response, nil
}

func (s *GrpcService) modelInfer(in *io.RecRequest, ServiceConfig *config_loader.ServiceConfig) (*RecommendResponse, error) {
	response := &RecommendResponse{
		Code: 404,
	}
	requestId := utils.CreateRequestId(in)

	//build model by model_factory
	modelName := in.GetModelType()
	if modelName != "" {
		modelName = strings.ToLower(modelName)
	}

	//strategy pattern. share model
	var modelStrategy model.ModelStrategyInterface
	modelStrategyContext := model.ModelStrategyContext{}
	_, ok := model.ShareModelsMap[in.GetDataId()]
	if !ok {
		modelfactory := model.ModelStrategyFactory{}
		modelStrategy = modelfactory.CreateModelStrategy(modelName, ServiceConfig)
		model.ShareModelsMap[in.GetDataId()] = modelStrategy
	} else {
		modelStrategy = model.ShareModelsMap[in.GetDataId()]
	}
	modelStrategyContext.SetModelStrategy(modelStrategy)

	var err error
	var responseTf map[string]interface{}

	if s.skywalkingWeatherOpen {
		responseTf, err = modelStrategyContext.ModelInferSkywalking(requestId, in.GetDataId(), in.GetItemList(), nil)
	} else {
		responseTf, err = modelStrategyContext.ModelInferNoSkywalking(requestId, in.GetDataId(), in.GetItemList(), nil)
	}
	if err != nil {
		logs.Error(requestId, time.Now(), err)
		return response, err
	}

	result := make([]*ItemInfo, 0)
	resultList := responseTf["data"].([]map[string]interface{})
	rankCh := make(chan *ItemInfo, len(resultList))
	for i := 0; i < len(resultList); i++ {
		inferWg.Add(1)
		go formatGrpcResponse(resultList[i], rankCh)
	}

	inferWg.Wait()
	close(rankCh)
	for itemScore := range rankCh {
		result = append(result, itemScore)
	}

	response = &RecommendResponse{
		Code:    200,
		Message: "success",
		Data: &ItemInfoList{
			Iteminfo_: result,
		},
	}

	return response, nil
}

func (s *GrpcService) modelInferReduce(in *io.RecRequest, ServiceConfig *config_loader.ServiceConfig) (*RecommendResponse, error) {
	response := &RecommendResponse{
		Code: 404,
	}
	requestId := utils.CreateRequestId(in)

	//build model by model_factory
	// modelName := in.GetModelType()
	// if modelName != "" {
	// 	modelName = strings.ToLower(modelName)
	// }

	modelName := "fm"

	//strategy pattern. share model
	var modelStrategy model.ModelStrategyInterface
	modelStrategyContext := model.ModelStrategyContext{}
	_, ok := model.ShareModelsMap[in.GetDataId()]
	if !ok {
		modelfactory := model.ModelStrategyFactory{}
		modelStrategy = modelfactory.CreateModelStrategy(modelName, ServiceConfig)
		model.ShareModelsMap[in.GetDataId()] = modelStrategy
	} else {
		modelStrategy = model.ShareModelsMap[in.GetDataId()]
	}
	modelStrategyContext.SetModelStrategy(modelStrategy)

	responseTf, err := modelStrategyContext.ModelInferSkywalking(requestId, in.GetDataId(), in.GetItemList(), nil)
	if err != nil {
		logs.Error(requestId, time.Now(), err)
		return response, err
	}

	result := make([]*ItemInfo, 0)
	resultList := responseTf["data"].([]map[string]interface{})
	rankCh := make(chan *ItemInfo, len(resultList))
	for i := 0; i < len(resultList); i++ {
		inferWg.Add(1)
		go formatGrpcResponse(resultList[i], rankCh)
	}

	inferWg.Wait()
	close(rankCh)
	for itemScore := range rankCh {
		result = append(result, itemScore)
	}

	response = &RecommendResponse{
		Code:    200,
		Message: "success",
		Data: &ItemInfoList{
			Iteminfo_: result,
		},
	}

	return response, nil
}

func formatGrpcResponse(itemScore map[string]interface{}, rankCh chan *ItemInfo) {
	defer inferWg.Done()

	itemid := itemScore["itemid"].(string)
	score := float32(itemScore["score"].(float64))
	itemInfo := &ItemInfo{
		Itemid: itemid,
		Score:  score,
	}

	rankCh <- itemInfo
}
