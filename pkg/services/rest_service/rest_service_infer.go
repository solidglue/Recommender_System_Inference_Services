package rest_service

import (
	"context"
	"errors"
	"infer-microservices/pkg/logs"
	"infer-microservices/pkg/utils"
	"net/http"
	"strings"

	config_loader "infer-microservices/pkg/config_loader"
	"infer-microservices/pkg/model"
	nacos "infer-microservices/pkg/nacos"
	"infer-microservices/pkg/services/io"
	"time"

	jsoniter "github.com/json-iterator/go"

	"github.com/afex/hystrix-go/hystrix"
)

type HttpService struct {
	nacosIp               string
	nacosPort             uint
	skywalkingWeatherOpen bool
	skywalkingIp          string
	skywalkingPort        uint
	skywalkingServerName  string
	lowerRankNum          int
	lowerRecallNum        int
	request               *http.Request
}

// set func

func (s *HttpService) SetNacosIp(nacosIp string) {
	s.nacosIp = nacosIp
}

func (s *HttpService) SetNacosPort(nacosPort uint) {
	s.nacosPort = nacosPort
}

func (s *HttpService) SetSkywalkingWeatherOpen(skywalkingWeatherOpen bool) {
	s.skywalkingWeatherOpen = skywalkingWeatherOpen
}

func (s *HttpService) GetSkywalkingWeatherOpen() bool {
	return s.skywalkingWeatherOpen
}

func (s *HttpService) SetSkywalkingIp(skywalkingIp string) {
	s.skywalkingIp = skywalkingIp
}

func (s *HttpService) GetSkywalkingIp() string {
	return s.skywalkingIp
}

func (s *HttpService) SetSkywalkingPort(skywalkingPort uint) {
	s.skywalkingPort = skywalkingPort
}

func (s *HttpService) GetSkywalkingPort() uint {
	return s.skywalkingPort
}

func (s *HttpService) SetSkywalkingServerName(skywalkingServerName string) {
	s.skywalkingServerName = skywalkingServerName
}

func (s *HttpService) GetSkywalkingServerName() string {
	return s.skywalkingServerName
}

func (s *HttpService) SetLowerRankNum(lowerRankNum int) {
	s.lowerRankNum = lowerRankNum
}

func (s *HttpService) SetLowerRecallNum(lowerRecallNum int) {
	s.lowerRecallNum = lowerRecallNum
}

// infer
func (s *HttpService) RecommenderInfer(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if info := recover(); info != nil {
			logs.Fatal("panic", info)
			rsp := make(map[string]interface{}, 0)
			rsp["error"] = "fatal"
			rsp["status"] = "fail"
			buff, _ := jsoniter.Marshal(rsp)
			w.Write(buff)

		} //else {
		//fmt.Println("")
		//}
	}()

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Millisecond*150)
	defer cancelFunc()
	r = r.WithContext(ctx)

	rsp := make(map[string]interface{}, 0)
	rsp["code"] = 200

	//INFO: convert http string data to struct data.
	request, err := s.convertHttpRequstToRecRequest(r)
	requestId := utils.CreateRequestId(&request)

	if err != nil {
		logs.Error(requestId, time.Now(), err)
		panic(err)
	}

	//check http input
	checkStatus := s.Check(requestId)
	if !checkStatus {
		err := errors.New("http input check failed")
		logs.Error(err)
		panic(err)
	}
	//check input
	checkStatus = request.Check()
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
	ServiceConfig := config_loader.ServiceConfigs[request.GetDataId()]
	response, err := s.recommenderInferHystrix("restServer", r, &request, ServiceConfig)
	if err != nil {
		logs.Error(requestId, time.Now(), err)
		panic(err)
	} else {
		rsp["data"] = response
	}

	buff, _ := jsoniter.Marshal(rsp)
	w.Write(buff)
}

func (s *HttpService) Check(requestId string) bool {
	err := s.request.ParseForm()
	if err != nil {
		logs.Error(requestId, time.Now(), err)
		return false
	}

	method := s.request.Method
	if method != "POST" {
		logs.Error(requestId, time.Now(), err)
		return false
	}

	data := s.request.Form["data"]
	if len(data) == 0 {
		logs.Error(requestId, time.Now(), err)
		return false
	}

	return true
}

func (s *HttpService) convertHttpRequstToRecRequest(r *http.Request) (io.RecRequest, error) {
	request := io.RecRequest{}
	data := s.request.Form["data"]
	requestMap := make(map[string]interface{}, 0)
	err := jsoniter.Unmarshal([]byte(data[0]), &requestMap)
	if err != nil {
		return request, err
	}

	return request, nil
}

func (s *HttpService) recommenderInferHystrix(serverName string, r *http.Request, in *io.RecRequest, ServiceConfig *config_loader.ServiceConfig) (map[string]interface{}, error) {
	response := make(map[string]interface{}, 0)
	requestId := utils.CreateRequestId(in)

	hystrixErr := hystrix.Do(serverName, func() error {
		// request recall / rank func.
		response_, err := s.modelInfer(r, in, ServiceConfig)
		if err != nil {
			logs.Error(requestId, time.Now(), err)
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
		response_, err_ := s.modelInferReduce(r, in, ServiceConfig)
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

func (s *HttpService) modelInfer(r *http.Request, in *io.RecRequest, ServiceConfig *config_loader.ServiceConfig) (map[string]interface{}, error) {
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

	var err error
	if s.skywalkingWeatherOpen {
		response, err = modelStrategyContext.ModelInferSkywalking(requestId, in.GetDataId(), in.GetItemList(), r)
	} else {
		response, err = modelStrategyContext.ModelInferNoSkywalking(requestId, in.GetDataId(), in.GetItemList(), r)
	}
	if err != nil {
		logs.Error(requestId, time.Now(), err)
		return response, err
	}

	return response, nil
}

func (s *HttpService) modelInferReduce(r *http.Request, in *io.RecRequest, ServiceConfig *config_loader.ServiceConfig) (map[string]interface{}, error) {
	response := make(map[string]interface{}, 0)
	requestId := utils.CreateRequestId(in)

	//build model by model_factory
	// modelName := in.GetModelType()
	// if modelName != "" {
	// 	modelName = strings.ToLower(modelName)
	// }

	modelName := "fm" //fm model use to reduce

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
	if s.skywalkingWeatherOpen {
		response, err = modelStrategyContext.ModelInferSkywalking(requestId, in.GetDataId(), in.GetItemList(), r)
	} else {
		response, err = modelStrategyContext.ModelInferNoSkywalking(requestId, in.GetDataId(), in.GetItemList(), r)
	}
	if err != nil {
		logs.Error(requestId, time.Now(), err)
		return response, err
	}

	return response, nil
}
