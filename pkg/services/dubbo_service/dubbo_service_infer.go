package server

import (
	"context"
	"errors"
	"fmt"
	config_loader "infer-microservices/pkg/config_loader"
	"infer-microservices/pkg/logs"
	"infer-microservices/pkg/model"
	nacos "infer-microservices/pkg/nacos"
	"infer-microservices/pkg/services/io"
	"infer-microservices/pkg/utils"
	"strings"
	"sync"
	"time"

	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"github.com/afex/hystrix-go/hystrix"
)

var recallWg sync.WaitGroup

type DubboService struct {
	nacosIp        string
	nacosPort      uint
	lowerRankNum   int
	lowerRecallNum int
}

//INFO:DONT REMOVE.  JAVA request service need it.
// // MethodMapper mapper upper func name to lower func name ,for java request.
// func (s *InferDubbogoService) MethodMapper() map[string]string {
// 	return map[string]string{
// 		"DubboRecommendServer": "dubboRecommendServer",
// 	}
// }

// set func
func (s *DubboService) SetNacosIp(nacosIp string) {
	s.nacosIp = nacosIp
}

func (s *DubboService) SetNacosPort(nacosPort uint) {
	s.nacosPort = nacosPort
}

func (s *DubboService) SetLowerRankNum(lowerRankNum int) {
	s.lowerRankNum = lowerRankNum
}

func (s *DubboService) SetLowerRecallNum(lowerRecallNum int) {
	s.lowerRecallNum = lowerRecallNum
}

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
	nacosConfig := nacosFactory.CreateNacosConfig(s.nacosIp, uint64(s.nacosPort), in)
	logs.Debug(requestId, time.Now(), "nacosConfig:", nacosConfig)

	nacosConfig.StartListenNacos()

	//infer
	ServiceConfig := config_loader.ServiceConfigs[in.GetDataId()]
	response_, err := s.recommenderInferHystrix("dubboServer", in, ServiceConfig)
	if err != nil {
		response.SetMessage(fmt.Sprintf("%s", err))
		panic(err)
	} else {
		response = response_
	}

	respCh <- response
}

func (s *DubboService) recommenderInferHystrix(serverName string, in *io.RecRequest, ServiceConfig *config_loader.ServiceConfig) (*io.RecResponse, error) {
	response := &io.RecResponse{}
	response.SetCode(404)
	requestId := utils.CreateRequestId(in)

	hystrixErr := hystrix.Do(serverName, func() error {
		// request recall / rank func.
		response_, err_ := s.modelInfer(in, ServiceConfig)
		if err_ != nil {
			logs.Error(requestId, time.Now(), err_)
			return err_
		} else {
			response = response_
		}
		logs.Debug(requestId, time.Now(), "hystrix unreduce response:", response)

		return err_
	}, func(err error) error {
		//INFO: do this when services are timeout (hystrix timeout).
		// less items and simple model.

		//INFO:its better not use the same func
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
		logs.Debug(requestId, time.Now(), "hystrix reduce response:", response)

		return err
	})

	if hystrixErr != nil {
		return response, hystrixErr
	}

	return response, nil
}

func (s *DubboService) modelInfer(in *io.RecRequest, ServiceConfig *config_loader.ServiceConfig) (*io.RecResponse, error) {
	response := &io.RecResponse{}
	response.SetCode(404)
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

	result, err := modelStrategyContext.ModelInferSkywalking(requestId, in.GetUserId(), in.GetItemList(), nil)
	if err != nil {
		logs.Error(requestId, time.Now(), err)
		return response, err
	}

	//package infer result.
	itemsScores := make([]string, 0)
	resultList := result["data"].([]map[string]interface{})
	recallCh := make(chan string, len(resultList))
	if len(resultList) > 0 {
		for i := 0; i < len(resultList); i++ {
			recallWg.Add(1)
			go formatDubboResponse(resultList[i], recallCh)
		}
		recallWg.Wait()
		close(recallCh)
		for itemScore := range recallCh {
			itemsScores = append(itemsScores, itemScore)
		}
		response.SetCode(200)
		response.SetMessage("success")
		response.SetData(itemsScores)
	}

	return response, nil
}

func (s *DubboService) modelInferReduce(in *io.RecRequest, ServiceConfig *config_loader.ServiceConfig) (*io.RecResponse, error) {
	response := &io.RecResponse{}
	response.SetCode(404)
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

	result, err := modelStrategyContext.ModelInferSkywalking(requestId, in.GetDataId(), in.GetItemList(), nil)
	if err != nil {
		logs.Error(requestId, time.Now(), err)
		return response, err
	}
	//package infer result.
	itemsScores := make([]string, 0)
	resultList := result["data"].([]map[string]interface{})
	recallCh := make(chan string, len(resultList))

	if len(resultList) > 0 {
		for i := 0; i < len(resultList); i++ {
			recallWg.Add(1)
			go formatDubboResponse(resultList[i], recallCh)
		}
		recallWg.Wait()
		close(recallCh)
		for itemScore := range recallCh {
			itemsScores = append(itemsScores, itemScore)
		}
		response.SetCode(200)
		response.SetMessage("success")
		response.SetData(itemsScores)
	}

	return response, nil
}

func formatDubboResponse(itemScore map[string]interface{}, recallCh chan string) {
	defer recallWg.Done()

	itemId := itemScore["itemid"].(string)
	score := float32(itemScore["score"].(float64))

	itemInfo := io.ItemInfo{}
	itemInfo.SetItemId(itemId)
	itemInfo.SetScore(score)

	itemScoreStr := utils.ConvertStructToJson(itemInfo)
	recallCh <- itemScoreStr
}
