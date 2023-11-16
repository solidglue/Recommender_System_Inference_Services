package grpc_service

import (
	"errors"
	"fmt"
	config_loader "infer-microservices/pkg/config_loader"
	"infer-microservices/pkg/model"
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

	ctx, cancelFunc := context.WithTimeout(ctx, time.Millisecond*150)
	defer cancelFunc()

	respCh := make(chan *RecommendResponse, 100)
	go g.recommenderInferContext(ctx, in, respCh)
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

func (s *GrpcService) recommenderInferContext(ctx context.Context, in *RecommendRequest, respCh chan *RecommendResponse) {
	defer func() {
		if info := recover(); info != nil {
			fmt.Println("panic", info)
		} //else {
		//fmt.Println("")
		//}
	}()

	response := &RecommendResponse{
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
	nacosFactory := nacos.NacosFactory{}
	nacosConfig := nacosFactory.CreateNacosConfig(s.nacosIp, uint64(s.nacosPort), &request)
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

	hystrixErr := hystrix.Do(serverName, func() error {
		// request recall / rank func.
		response_, err := s.modelInfer(in, ServiceConfig)
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
		response_, err_ := s.modelInferReduce(in, ServiceConfig)

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

func (s *GrpcService) modelInfer(in *io.RecRequest, ServiceConfig *config_loader.ServiceConfig) (*RecommendResponse, error) {
	response := &RecommendResponse{
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
