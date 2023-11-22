package baseservice

import (
	config_loader "infer-microservices/pkg/config_loader"
	"infer-microservices/pkg/logs"
	"infer-microservices/pkg/model"
	"infer-microservices/pkg/model/basemodel"
	"infer-microservices/pkg/services/io"
	"infer-microservices/pkg/utils"
	"net/http"
	"strings"
	"time"

	"github.com/afex/hystrix-go/hystrix"
)

// baseservice, all service extend baseservice
type BaseService struct {
	nacosIp               string
	nacosPort             uint
	skywalkingWeatherOpen bool
	skywalkingIp          string
	skywalkingPort        uint
	skywalkingServerName  string
	lowerRankNum          int
	lowerRecallNum        int
}

// set func
func (s *BaseService) SetNacosIp(nacosIp string) {
	s.nacosIp = nacosIp
}

func (s *BaseService) GetNacosIp() string {
	return s.nacosIp
}

func (s *BaseService) SetNacosPort(nacosPort uint) {
	s.nacosPort = nacosPort
}

func (s *BaseService) GetNacosPort() uint {
	return s.nacosPort
}

func (s *BaseService) SetSkywalkingWeatherOpen(skywalkingWeatherOpen bool) {
	s.skywalkingWeatherOpen = skywalkingWeatherOpen
}

func (s *BaseService) GetSkywalkingWeatherOpen() bool {
	return s.skywalkingWeatherOpen
}

func (s *BaseService) SetSkywalkingIp(skywalkingIp string) {
	s.skywalkingIp = skywalkingIp
}

func (s *BaseService) GetSkywalkingIp() string {
	return s.skywalkingIp
}

func (s *BaseService) SetSkywalkingPort(skywalkingPort uint) {
	s.skywalkingPort = skywalkingPort
}

func (s *BaseService) GetSkywalkingPort() uint {
	return s.skywalkingPort
}

func (s *BaseService) SetSkywalkingServerName(skywalkingServerName string) {
	s.skywalkingServerName = skywalkingServerName
}

func (s *BaseService) GetSkywalkingServerName() string {
	return s.skywalkingServerName
}

func (s *BaseService) SetLowerRankNum(lowerRankNum int) {
	s.lowerRankNum = lowerRankNum
}

func (s *BaseService) GetLowerRankNum() int {
	return s.lowerRankNum
}

func (s *BaseService) SetLowerRecallNum(lowerRecallNum int) {
	s.lowerRecallNum = lowerRecallNum
}

func (s *BaseService) GetLowerRecallNum() int {
	return s.lowerRecallNum
}

func (s *BaseService) RecommenderInferHystrix(r *http.Request, serverName string, in *io.RecRequest, ServiceConfig *config_loader.ServiceConfig) (map[string]interface{}, error) {
	response := make(map[string]interface{}, 0)
	requestId := utils.CreateRequestId(in)

	hystrixErr := hystrix.Do(serverName, func() error {
		// request recall / rank func.
		response_, err_ := s.modelInfer(r, in, ServiceConfig)
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
		response_, err_ := s.modelInferReduce(r, in, ServiceConfig)
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

func (s *BaseService) modelInfer(r *http.Request, in *io.RecRequest, ServiceConfig *config_loader.ServiceConfig) (map[string]interface{}, error) {
	response := make(map[string]interface{}, 0)
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

	//use callback func to create sample
	createSampleFunc := basemodel.SampleCallBackFuncMap[strings.ToLower(modelStrategy.GetModelType())]
	result, err := modelStrategyContext.ModelInferSkywalking(requestId, in.GetUserId(), in.GetItemList(), r, createSampleFunc)
	if err != nil {
		logs.Error(requestId, time.Now(), err)
		return response, err
	}

	//package infer result.
	itemsScores := make([]*io.ItemInfo, 0)
	resultList := result["data"].([]map[string]interface{})
	resultCh := make(chan *io.ItemInfo, len(resultList))
	if len(resultList) > 0 {
		for i := 0; i < len(resultList); i++ {
			go formatDubboResponse(resultList[i], resultCh)
		}

	loop:
		for {
			select {
			case <-time.After(time.Millisecond * 100):
				break loop
			case itemScore := <-resultCh:
				itemsScores = append(itemsScores, itemScore)
			}
		}
		close(resultCh)

		response["code"] = 200
		response["message"] = "success"
		response["data"] = itemsScores
	}

	return response, nil
}

func (s *BaseService) modelInferReduce(r *http.Request, in *io.RecRequest, ServiceConfig *config_loader.ServiceConfig) (map[string]interface{}, error) {
	response := make(map[string]interface{}, 0)
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

	//use callback func to create sample
	createSampleFunc := basemodel.SampleCallBackFuncMap[strings.ToLower(modelStrategy.GetModelType())]
	result, err := modelStrategyContext.ModelInferSkywalking(requestId, in.GetDataId(), in.GetItemList(), r, createSampleFunc)
	if err != nil {
		logs.Error(requestId, time.Now(), err)
		return response, err
	}
	//package infer result.
	itemsScores := make([]*io.ItemInfo, 0)
	resultList := result["data"].([]map[string]interface{})
	resultCh := make(chan *io.ItemInfo, len(resultList))

	if len(resultList) > 0 {
		for i := 0; i < len(resultList); i++ {
			go formatDubboResponse(resultList[i], resultCh)
		}

	loop:
		for {
			select {
			case <-time.After(time.Millisecond * 100):
				break loop
			case itemScore := <-resultCh:
				itemsScores = append(itemsScores, itemScore)
			}
		}
		close(resultCh)

		response["code"] = 200
		response["message"] = "success"
		response["data"] = itemsScores
	}

	return response, nil
}

func formatDubboResponse(itemScore map[string]interface{}, resultCh chan *io.ItemInfo) { //recallCh chan string) {
	itemId := itemScore["itemid"].(string)
	score := float32(itemScore["score"].(float64))

	itemInfo := io.ItemInfo{}
	itemInfo.SetItemId(itemId)
	itemInfo.SetScore(score)

	//itemScoreStr := utils.ConvertStructToJson(itemInfo)
	//recallCh <- itemScoreStr

	resultCh <- &itemInfo
}
