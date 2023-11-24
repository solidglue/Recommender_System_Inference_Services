package rest_service

import (
	"errors"
	config_loader "infer-microservices/pkg/config_loader"
	"infer-microservices/pkg/logs"
	nacos "infer-microservices/pkg/nacos"
	"infer-microservices/pkg/services/baseservice"
	"infer-microservices/pkg/services/io"
	"infer-microservices/pkg/utils"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo"
)

// extend from  baseservice
type EchoService struct {
	baseservice *baseservice.BaseService
}

func (s *EchoService) SetBaseService(baseservice *baseservice.BaseService) {
	s.baseservice = baseservice
}

func (s *EchoService) GetBaseService() *baseservice.BaseService {
	return s.baseservice
}

// sync server
func (s *EchoService) SyncRecommenderInfer(c echo.Context) error {
	respCh := make(chan map[string]interface{}, 100)
	go s.RecommenderInfer(c, respCh)

	select {
	case <-time.After(time.Millisecond * 100):
		rsp := make(map[string]interface{}, 0)
		return c.JSON(http.StatusRequestTimeout, rsp)
	case responseCh := <-respCh:
		return c.JSON(http.StatusOK, responseCh)
	}

}

// infer
func (s *EchoService) RecommenderInfer(c echo.Context, ch chan<- map[string]interface{}) {
	defer func() {
		if info := recover(); info != nil {
			logs.Fatal("panic", info)
			rsp := make(map[string]interface{}, 0)
			ch <- rsp
		} //else {
		//fmt.Println("")
		//}
	}()

	rsp := make(map[string]interface{}, 0)
	//INFO: convert http string data to struct data.
	request, err := s.convertHttpRequstToRecRequest(c)
	requestId := utils.CreateRequestId(&request)

	if err != nil {
		logs.Error(requestId, time.Now(), err)
		panic(err)
	}

	//check http input
	checkStatus := s.check(c, requestId)
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
	response, err := s.baseservice.RecommenderInferHystrix(c.Request(), "restServer", &request, ServiceConfig)
	if err != nil {
		logs.Error(requestId, time.Now(), err)
		panic(err)
	} else {
		rsp["data"] = response
	}

	ch <- rsp
}

func (s *EchoService) check(c echo.Context, requestId string) bool {
	err := c.Request().ParseForm()
	if err != nil {
		logs.Error(requestId, time.Now(), err)
		return false
	}

	method := c.Request().Method
	if method != "POST" {
		logs.Error(requestId, time.Now(), err)
		return false
	}

	data := c.Request().Form["data"]
	if len(data) == 0 {
		logs.Error(requestId, time.Now(), err)
		return false
	}

	return true
}

func (s *EchoService) convertHttpRequstToRecRequest(c echo.Context) (io.RecRequest, error) {
	request := io.RecRequest{}
	data := c.Request().Form["data"]
	requestMap := make(map[string]interface{}, 0)
	err := jsoniter.Unmarshal([]byte(data[0]), &requestMap)
	if err != nil {
		return request, err
	}

	return request, nil
}
