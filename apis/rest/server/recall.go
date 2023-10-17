package server

import (
	"context"
	"errors"
	"infer-microservices/apis"
	"infer-microservices/apis/input_format"
	"infer-microservices/common/flags"
	"infer-microservices/cores/nacos_config_listener"
	"infer-microservices/cores/service_config_loader"
	"infer-microservices/utils/logs"
	"net/http"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	jsoniter "github.com/json-iterator/go"
)

var lowerRecallNum int

type recallServer struct {
}

func init() {
	flagFactory := flags.FlagFactory{}
	flagHystrix := flagFactory.FlagHystrixFactory()

	lowerRecallNum = *flagHystrix.GetHystrixLowerRecallNum()
	lowerRankNum = *flagHystrix.GetHystrixLowerRankNum()
}

// emb
func (c *recallServer) restInferServer(w http.ResponseWriter, r *http.Request) {

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

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Millisecond*100)
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
	}

	//INFO: convert http string data to struct data.
	request, err := httpRequstParse(r)
	if err != nil {
		logs.Error(err)
		panic(err)
	}

	ServiceConfig := apis.ServiceConfigs[request.GetDataId()]
	response, err := c.restHystrixRecaller("restServer", r, &request, ServiceConfig)
	if err != nil {
		logs.Error(err)
		panic(err)
	} else {
		rsp["data"] = response
	}

	buff, _ := jsoniter.Marshal(rsp)
	w.Write(buff)
}

func (c *recallServer) restHystrixRecaller(serverName string, r *http.Request, in *apis.RecRequest, ServiceConfig *service_config_loader.ServiceConfig) (map[string]interface{}, error) {
	response := make(map[string]interface{}, 0)
	hystrixErr := hystrix.Do(serverName, func() error {
		// request recall / rank func.
		response_, err := c.restRecaller(r, in, ServiceConfig)
		if err != nil {
			logs.Error(err)
		} else {
			response = response_
		}

		return nil
	}, func(err error) error {
		// do this when services are timeout
		if err != nil {
			logs.Error(err)
		}

		itemList := in.GetItemList()
		in.SetRecallNum(int32(lowerRecallNum))
		in.SetItemList(itemList[:lowerRankNum])
		response_, err := c.restRecaller(r, in, ServiceConfig)
		if err != nil {
			logs.Error(err)
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

func (c *recallServer) restRecaller(r *http.Request, in *apis.RecRequest, ServiceConfig *service_config_loader.ServiceConfig) (map[string]interface{}, error) {
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

	_, ok := apis.NacosListedMap[dataId]
	if !ok {
		err := nacosConn.ServiceConfigListen()
		if err != nil {
			return response, err
		} else {
			apis.NacosListedMap[dataId] = true
		}
	}

	recaller := input_format.RecallInputFormat{}
	dssm, err := recaller.InputCheckAndFormat(in, ServiceConfig)
	if err != nil {
		logs.Error(err)
		return response, err
	}

	if skywalkingWeatherOpen {
		response, err = dssm.RecallInferSkywalking(r)
	} else {
		response, err = dssm.RecallInferNoSkywalking(r)
	}

	if err != nil {
		logs.Error(err)
		return response, err
	}

	return response, nil
}
