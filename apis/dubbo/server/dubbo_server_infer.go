package server

import (
	"context"
	"errors"
	"fmt"
	"infer-microservices/apis/io"
	"infer-microservices/cores/model"
	"infer-microservices/cores/nacos_config_listener"
	"infer-microservices/cores/service_config_loader"
	"infer-microservices/utils"
	"infer-microservices/utils/logs"
	"strings"
	"sync"
	"time"

	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"github.com/afex/hystrix-go/hystrix"
)

// var inferModel model.ModelInferInterface
var recallWg sync.WaitGroup

//INFO:DONT REMOVE.  JAVA request service need it.
// // MethodMapper mapper upper func name to lower func name ,for java request.
// func (s *InferDubbogoService) MethodMapper() map[string]string {
// 	return map[string]string{
// 		"DubboRecommendServer": "dubboRecommendServer",
// 	}
// }

// Implement interface methods.
func (s *DubboServer) DubboRecommendServer(ctx context.Context, in *io.RecRequest) (*io.RecResponse, error) {
	response := &io.RecResponse{}
	response.SetCode(404)

	//check input
	checkStatus := in.Check()
	if !checkStatus {
		err := errors.New("input check failed")
		logs.Error(err)
		return response, err
	}

	//INFO: set timeout by context, degraded service by hystix.
	ctx, cancelFunc := context.WithTimeout(ctx, time.Millisecond*100)
	defer cancelFunc()

	respCh := make(chan *io.RecResponse, 100)
	go s.dubboRecommenderServerContext(ctx, in, respCh)

	for {
		select {
		case <-ctx.Done():
			switch ctx.Err() {
			case context.DeadlineExceeded:
				logs.Info("context timeout DeadlineExceeded.")
				return response, ctx.Err()
			case context.Canceled:
				logs.Info("context timeout Canceled.")
				return response, ctx.Err()
			}
		case responseCh := <-respCh:
			response = responseCh
			return response, nil
		}
	}
}

func (s *DubboServer) dubboRecommenderServerContext(ctx context.Context, in *io.RecRequest, respCh chan *io.RecResponse) {
	defer func() {
		if info := recover(); info != nil {
			fmt.Println("panic", info)
		} //else {
		//  fmt.Println("finish.")
		//}
	}()

	response := &io.RecResponse{}
	response.SetCode(404)

	//nacos listen
	nacosFactory := nacos_config_listener.NacosFactory{}
	nacosConfig := nacosFactory.CreateNacosConfig(s.nacosIp, uint64(s.nacosPort), in)
	nacosConfig.StartListenNacos()

	//infer
	ServiceConfig := service_config_loader.ServiceConfigs[in.GetDataId()]
	response_, err := s.dubboHystrixServer("dubboServer", in, ServiceConfig)
	if err != nil {
		response.SetMessage(fmt.Sprintf("%s", err))
		panic(err)
	} else {
		response = response_
	}

	respCh <- response
}

func (s *DubboServer) dubboHystrixServer(serverName string, in *io.RecRequest, ServiceConfig *service_config_loader.ServiceConfig) (*io.RecResponse, error) {
	response := &io.RecResponse{}
	response.SetCode(404)

	hystrixErr := hystrix.Do(serverName, func() error {
		// request recall / rank func.
		response_, err_ := s.recommenderInfer(in, ServiceConfig)
		if err_ != nil {
			logs.Error(err_)
			return err_
		} else {
			response = response_
		}
		return nil
	}, func(err error) error {
		//INFO: do this when services are timeout (hystrix timeout).
		// less items and simple model.

		//INFO:its better not use the same func
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

func (s *DubboServer) recommenderInfer(in *io.RecRequest, ServiceConfig *service_config_loader.ServiceConfig) (*io.RecResponse, error) {
	response := &io.RecResponse{}
	response.SetCode(404)

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

	result, err := modelStrategyContext.ModelInferSkywalking(nil)
	if err != nil {
		logs.Error(err)
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

func (s *DubboServer) recommenderInferReduce(in *io.RecRequest, ServiceConfig *service_config_loader.ServiceConfig) (*io.RecResponse, error) {
	response := &io.RecResponse{}
	response.SetCode(404)

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
	result, err := modelStrategyContext.ModelInferSkywalking(nil)
	if err != nil {
		logs.Error(err)
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
