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

	err := r.ParseForm()
	if err != nil {
		rsp["code"] = 404
		rsp["error"] = errors.New("ParseForm Error")
		panic(err)
	}

	method := r.Method
	if method != "POST" {
		rsp["code"] = 404
		rsp["error"] = errors.New("method should be POST")
		panic(err)

	}

	data := r.Form["data"]
	if len(data) == 0 {
		rsp["code"] = 404
		rsp["error"] = errors.New("emt input data")
		panic(err)
	}

	//INFO: convert http string data to struct data.
	request, err := httpRequstParse(r)
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

func httpRequstParse(r *http.Request) (io.RecRequest, error) {
	request := io.RecRequest{}

	err := r.ParseForm()
	if err != nil {
		return request, err
	}

	method := r.Method
	if method != "POST" {
		return request, err
	}

	data := r.Form["data"]
	if len(data) == 0 {
		return request, err
	}

	requestMap := make(map[string]interface{}, 0)
	err = jsoniter.Unmarshal([]byte(data[0]), &requestMap)
	if err != nil {
		return request, err
	}

	request, err = inputCheck(requestMap)
	if err != nil {
		return request, err
	}

	return request, nil
}

func inputCheck(requestMap map[string]interface{}) (io.RecRequest, error) {
	request := io.RecRequest{}

	//dataId
	dataId, ok := requestMap["dataId"]
	if ok {
		request.SetDataId(dataId.(string))
	} else {
		return request, errors.New("dataId can not be empty")
	}

	//modelType
	modelType, ok := requestMap["modelType"]
	if ok {
		request.SetModelType(modelType.(string))
	} else {
		return request, errors.New("modelType can not be empty")
	}

	//userId
	userId, ok := requestMap["userId"]
	if ok {
		request.SetUserId(userId.(string))
	} else {
		return request, errors.New("userId can not be empty")
	}

	//recall num. reflect.
	recallNum := int32(100)
	recallNumType := reflect.TypeOf(requestMap["recallNum"])
	recallNumTypeKind := recallNumType.Kind()
	switch recallNumTypeKind {
	case reflect.String:
		recallNumStr, ok0 := requestMap["recallNum"].(string)
		if ok0 {
			recallNum64, err := strconv.ParseInt(recallNumStr, 10, 64)
			if err != nil {
				ok = false
			} else {
				recallNum = int32(recallNum64)
				ok = true
			}
		}
	case reflect.Float32, reflect.Float64, reflect.Int16, reflect.Int, reflect.Int64, reflect.Int8,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		recallNum, ok = requestMap["recallNum"].(int32)
	default:
		err := errors.New("unkown type, set recallnum to 100")
		logs.Error(err)
	}

	if ok {
		request.SetRecallNum(recallNum)
	} else {
		return request, errors.New("dataId can not be empty")
	}

	if recallNum > 1000 {
		return request, errors.New("recallNum should less than 1000 ")
	}

	//itemList
	itemList, ok := requestMap["itemList"].([]string)
	if ok {
		request.SetItemList(itemList)
	} else {
		return request, errors.New("itemList can not be empty")
	}

	if len(itemList) > 200 {
		return request, errors.New("itemList's len should less than 200 ")
	}

	return request, nil
}

func (s *HttpServer) restHystrixInfer(serverName string, r *http.Request, in *io.RecRequest, ServiceConfig *service_config_loader.ServiceConfig) (map[string]interface{}, error) {
	response := make(map[string]interface{}, 0)
	hystrixErr := hystrix.Do(serverName, func() error {
		// request recall / rank func.
		response_, err := s.recommenderInferReduce(r, in, ServiceConfig)
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

	modelfactory := model.ModelFactory{}
	modelinfer, err := modelfactory.CreateInferModel(modelName, in, ServiceConfig)
	if err != nil {
		return response, err
	}

	if s.skywalkingWeatherOpen {
		response, err = modelinfer.ModelInferSkywalking(r)
	} else {
		response, err = modelinfer.ModelInferNoSkywalking(r)
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

	modelfactory := model.ModelFactory{}
	modelinfer, err := modelfactory.CreateInferModel(modelName, in, ServiceConfig)
	if err != nil {
		return response, err
	}

	if s.skywalkingWeatherOpen {
		response, err = modelinfer.ModelInferSkywalking(r)
	} else {
		response, err = modelinfer.ModelInferNoSkywalking(r)
	}
	if err != nil {
		logs.Error(err)
		return response, err
	}

	return response, nil
}
