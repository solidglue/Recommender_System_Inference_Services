package rest_service

import (
	"context"
	"errors"
	"infer-microservices/internal/logs"
	"infer-microservices/internal/utils"
	"net/http"

	config_loader "infer-microservices/pkg/config_loader"
	nacos "infer-microservices/pkg/nacos"
	"infer-microservices/pkg/services/baseservice"
	"infer-microservices/pkg/services/io"
	"time"

	jsoniter "github.com/json-iterator/go"
)

// extend from  baseservice
type HttpService struct {
	baseservice *baseservice.BaseService
	request     *http.Request
}

func (s *HttpService) SetBaseService(baseservice *baseservice.BaseService) {
	s.baseservice = baseservice
}

func (s *HttpService) GetBaseService() *baseservice.BaseService {
	return s.baseservice
}

func (s *HttpService) SetRequest(request *http.Request) {
	s.request = request
}

func (s *HttpService) GetRequest() *http.Request {
	return s.request
}

// sync server
func (s *HttpService) SyncRecommenderInfer(w http.ResponseWriter, r *http.Request) {
	respCh := make(chan []byte, 100)
	go s.RecommenderInfer(w, r, respCh)

	select {
	case <-time.After(time.Millisecond * 100):
		rsp := make(map[string]interface{}, 0)
		rsp["error"] = "fatal"
		rsp["status"] = "fail"
		buff, _ := jsoniter.Marshal(rsp)
		w.Write(buff)
	case responseCh := <-respCh:
		w.Write(responseCh)
	}

}

// infer
func (s *HttpService) RecommenderInfer(w http.ResponseWriter, r *http.Request, ch chan<- []byte) {
	defer func() {
		if info := recover(); info != nil {
			logs.Fatal("panic", info)
			rsp := make(map[string]interface{}, 0)
			rsp["error"] = "fatal"
			rsp["status"] = "fail"
			buff, _ := jsoniter.Marshal(rsp)
			//w.Write(buff)
			ch <- buff
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
	nacosConfig := nacosFactory.CreateNacosConfig(s.baseservice.GetNacosIp(), uint64(s.baseservice.GetNacosPort()), &request)
	logs.Debug(requestId, time.Now(), "nacosConfig:", nacosConfig)

	nacosConfig.StartListenNacos()

	//infer
	ServiceConfig := config_loader.GetServiceConfigs()[request.GetDataId()]
	response, err := s.baseservice.RecommenderInferHystrix(r, "restServer", &request, ServiceConfig)
	if err != nil {
		logs.Error(requestId, time.Now(), err)
		panic(err)
	} else {
		rsp["data"] = response
	}

	buff, _ := jsoniter.Marshal(rsp)
	//w.Write(buff)

	ch <- buff
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
