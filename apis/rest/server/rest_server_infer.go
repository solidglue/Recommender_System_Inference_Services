package server

import (
	"context"
	"errors"
	"infer-microservices/utils/logs"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"infer-microservices/apis/io"
	"infer-microservices/cores/model"
	"infer-microservices/cores/nacos_config_listener"
	"infer-microservices/cores/service_config_loader"
	"time"

	jsoniter "github.com/json-iterator/go"

	"github.com/afex/hystrix-go/hystrix"
)

// infer
func (s *HttpServer) restInferServer(w http.ResponseWriter, r *http.Request) {
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

	//check http input
	checkStatus := s.Check()
	if !checkStatus {
		err := errors.New("http input check failed")
		logs.Error(err)
		panic(err)
	}

	//INFO: convert http string data to struct data.
	request, err := s.httpRequstParse(r)

	//check input
	checkStatus = request.Check()
	if !checkStatus {
		err := errors.New("input check failed")
		logs.Error(err)
		panic(err)
	}

	if err != nil {
		logs.Error(err)
		panic(err)
	}

	ServiceConfig := service_config_loader.ServiceConfigs[request.GetDataId()]
	response, err := s.restHystrixInfer("restServer", r, &request, ServiceConfig)
	if err != nil {
		logs.Error(err)
		panic(err)
	} else {
		rsp["data"] = response
	}

	buff, _ := jsoniter.Marshal(rsp)
	w.Write(buff)
}

func (s *HttpServer) Check() bool {
	err := s.request.ParseForm()
	if err != nil {
		logs.Error(err)
		return false
	}

	method := s.request.Method
	if method != "POST" {
		logs.Error(err)
		return false
	}

	data := s.request.Form["data"]
	if len(data) == 0 {
		logs.Error(err)
		return false
	}

	return true
}

func (s *HttpServer) httpRequstParse(r *http.Request) (io.RecRequest, error) {
	request := io.RecRequest{}
	data := s.request.Form["data"]
	requestMap := make(map[string]interface{}, 0)
	err := jsoniter.Unmarshal([]byte(data[0]), &requestMap)
	if err != nil {
		return request, err
	}

	//INFO:the recallNum param from http request,maybe int/ float /string。 user reflect to convert to int32.
	recallNum := int32(100)
	recallNumType := reflect.TypeOf(requestMap["recallNum"])
	recallNumTypeKind := recallNumType.Kind()
	switch recallNumTypeKind {
	case reflect.String:
		recallNumStr, ok := requestMap["recallNum"].(string)
		if ok {
			recallNum64, err := strconv.ParseInt(recallNumStr, 10, 64)
			if err != nil {
				logs.Error(err)
			} else {
				recallNum = int32(recallNum64)
			}
		}
	case reflect.Float32, reflect.Float64, reflect.Int16, reflect.Int, reflect.Int64, reflect.Int8,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		recallNum, _ = requestMap["recallNum"].(int32)
	default:
		logs.Info("unkown type, set recallnum to 100")
	}
	request.SetRecallNum(recallNum)

	return request, nil
}

func (s *HttpServer) restHystrixInfer(serverName string, r *http.Request, in *io.RecRequest, ServiceConfig *service_config_loader.ServiceConfig) (map[string]interface{}, error) {
	response := make(map[string]interface{}, 0)
	hystrixErr := hystrix.Do(serverName, func() error {
		// request recall / rank func.
		response_, err := s.recommenderInfer(r, in, ServiceConfig)
		if err != nil {
			logs.Error(err)
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
		response_, err_ := s.recommenderInferReduce(r, in, ServiceConfig)
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

func (s *HttpServer) recommenderInfer(r *http.Request, in *io.RecRequest, ServiceConfig *service_config_loader.ServiceConfig) (map[string]interface{}, error) {
	response := make(map[string]interface{}, 0)

	dataId := in.GetDataId()
	groupId := in.GetGroupId()
	namespaceId := in.GetNamespaceId()

	nacosConn := nacos_config_listener.NacosConnConfig{}
	nacosConn.SetDataId(dataId)
	nacosConn.SetGroupId(groupId)
	nacosConn.SetNamespaceId(namespaceId)
	nacosConn.SetIp(s.nacosIp)
	nacosConn.SetPort(uint64(s.nacosPort))

	_, ok := nacos_config_listener.NacosListedMap[dataId]
	if !ok {
		err := nacosConn.ServiceConfigListen()
		if err != nil {
			return response, err
		} else {
			nacos_config_listener.NacosListedMap[dataId] = true
		}
	}

	//build model by model_factory
	modelName := in.GetModelType()
	if modelName != "" {
		modelName = strings.ToLower(modelName)
	}

	//strategy pattern
	modelfactory := model.ModelStrategyFactory{}
	modelStrategyContext := model.ModelStrategyContext{}
	modelStrategy := modelfactory.CreateModelStrategy(modelName, in, ServiceConfig)
	modelStrategyContext.SetModelStrategy(modelStrategy)

	var err error
	if s.skywalkingWeatherOpen {
		response, err = modelStrategyContext.ModelInferSkywalking(r)
	} else {
		response, err = modelStrategyContext.ModelInferNoSkywalking(r)
	}
	if err != nil {
		logs.Error(err)
		return response, err
	}

	return response, nil
}

func (s *HttpServer) recommenderInferReduce(r *http.Request, in *io.RecRequest, ServiceConfig *service_config_loader.ServiceConfig) (map[string]interface{}, error) {
	response := make(map[string]interface{}, 0)

	dataId := in.GetDataId()
	groupId := in.GetGroupId()
	namespaceId := in.GetNamespaceId()

	nacosConn := nacos_config_listener.NacosConnConfig{}
	nacosConn.SetDataId(dataId)
	nacosConn.SetGroupId(groupId)
	nacosConn.SetNamespaceId(namespaceId)
	nacosConn.SetIp(s.nacosIp)
	nacosConn.SetPort(uint64(s.nacosPort))

	_, ok := nacos_config_listener.NacosListedMap[dataId]
	if !ok {
		err := nacosConn.ServiceConfigListen()
		if err != nil {
			return response, err
		} else {
			nacos_config_listener.NacosListedMap[dataId] = true
		}
	}

	//build model by model_factory
	// modelName := in.GetModelType()
	// if modelName != "" {
	// 	modelName = strings.ToLower(modelName)
	// }

	modelName := "fm" //fm model use to reduce

	//strategy pattern
	modelfactory := model.ModelStrategyFactory{}
	modelStrategy := modelfactory.CreateModelStrategy(modelName, in, ServiceConfig)
	modelStrategyContext := model.ModelStrategyContext{}
	modelStrategyContext.SetModelStrategy(modelStrategy)

	var err error
	if s.skywalkingWeatherOpen {
		response, err = modelStrategyContext.ModelInferSkywalking(r)
	} else {
		response, err = modelStrategyContext.ModelInferNoSkywalking(r)
	}
	if err != nil {
		logs.Error(err)
		return response, err
	}

	return response, nil
}
