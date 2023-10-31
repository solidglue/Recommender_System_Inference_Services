package server

import (
	"context"
	"errors"
	"strings"

	"infer-microservices/apis/io"
	"infer-microservices/common/flags"
	"infer-microservices/cores/model"
	"infer-microservices/cores/nacos_config_listener"
	"infer-microservices/cores/service_config_loader"
	"infer-microservices/utils/logs"
	"net/http"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	jsoniter "github.com/json-iterator/go"
)

var lowerRankNum int
var lowerRecallNum int

type RestInferService struct {
}

func init() {
	flagFactory := flags.FlagFactory{}
	flagHystrix := flagFactory.CreateFlagHystrix()

	lowerRecallNum = *flagHystrix.GetHystrixLowerRecallNum()
	lowerRankNum = *flagHystrix.GetHystrixLowerRankNum()
}

// infer
func (s *RestInferService) restInferServer(w http.ResponseWriter, r *http.Request) {

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

func (s *RestInferService) restHystrixInfer(serverName string, r *http.Request, in *io.RecRequest, ServiceConfig *service_config_loader.ServiceConfig) (map[string]interface{}, error) {
	response := make(map[string]interface{}, 0)
	hystrixErr := hystrix.Do(serverName, func() error {
		// request recall / rank func.
		response_, err := s.RecommenderInferReduce(r, in, ServiceConfig)
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
		in.SetRecallNum(int32(lowerRecallNum))
		in.SetItemList(itemList[:lowerRankNum])
		response_, err_ := s.RecommenderInferReduce(r, in, ServiceConfig)
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

func (s *RestInferService) RecommenderInfer(r *http.Request, in *io.RecRequest, ServiceConfig *service_config_loader.ServiceConfig) (map[string]interface{}, error) {
	response := make(map[string]interface{}, 0)

	dataId := in.GetDataId()
	groupId := in.GetGroupId()
	namespaceId := in.GetNamespaceId()

	nacosConn := nacos_config_listener.NacosConnConfig{}
	nacosConn.SetDataId(dataId)
	nacosConn.SetGroupId(groupId)
	nacosConn.SetNamespaceId(namespaceId)
	nacosConn.SetIp(NacosIP)
	nacosConn.SetPort(NacosPort)

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

	if skywalkingWeatherOpen {
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

func (s *RestInferService) RecommenderInferReduce(r *http.Request, in *io.RecRequest, ServiceConfig *service_config_loader.ServiceConfig) (map[string]interface{}, error) {
	response := make(map[string]interface{}, 0)

	dataId := in.GetDataId()
	groupId := in.GetGroupId()
	namespaceId := in.GetNamespaceId()

	nacosConn := nacos_config_listener.NacosConnConfig{}
	nacosConn.SetDataId(dataId)
	nacosConn.SetGroupId(groupId)
	nacosConn.SetNamespaceId(namespaceId)
	nacosConn.SetIp(NacosIP)
	nacosConn.SetPort(NacosPort)

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

	if skywalkingWeatherOpen {
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
