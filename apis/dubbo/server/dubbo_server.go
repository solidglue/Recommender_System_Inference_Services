package server

import (
	"context"
	"fmt"
	"infer-microservices/apis/io"
	"infer-microservices/common/flags"
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

var ipAddr_ string
var port_ uint64
var lowerRecallNum int
var lowerRankNum int

// var inferModel model.ModelInferInterface
var recallWg sync.WaitGroup

type DubbogoInferService struct {
}

func init() {
	// //regisger dubbo service.
	// hessian.RegisterPOJO(&io.RecRequest{})
	// hessian.RegisterPOJO(&io.RecResponse{})
	// config.SetProviderService(&DubbogoInferService{})

	//get hystrix parms.
	flagFactory := flags.FlagFactory{}
	flagHystrix := flagFactory.CreateFlagHystrix()
	lowerRecallNum = *flagHystrix.GetHystrixLowerRecallNum()
	lowerRankNum = *flagHystrix.GetHystrixLowerRankNum()
}

//INFO:DONT REMOVE.  JAVA request service need it.
// // MethodMapper mapper upper func name to lower func name ,for java request.
// func (s *InferDubbogoService) MethodMapper() map[string]string {
// 	return map[string]string{
// 		"DubboRecommendServer": "dubboRecommendServer",
// 	}
// }

// Implement interface methods.
func (r *DubbogoInferService) DubboRecommendServer(ctx context.Context, in *io.RecRequest) (*io.RecResponse, error) {
	response := &io.RecResponse{}
	response.SetCode(404)

	//INFO: set timeout by context, degraded service by hystix.
	ctx, cancelFunc := context.WithTimeout(ctx, time.Millisecond*100)
	defer cancelFunc()

	respCh := make(chan *io.RecResponse, 100)
	go r.dubboRecommenderServerContext(ctx, in, respCh)

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

func (r *DubbogoInferService) dubboRecommenderServerContext(ctx context.Context, in *io.RecRequest, respCh chan *io.RecResponse) {
	defer func() {
		if info := recover(); info != nil {
			fmt.Println("panic", info)
		} //else {
		//  fmt.Println("finish.")
		//}
	}()

	response := &io.RecResponse{}
	response.SetCode(404)

	nacosConn := getNacosConn(in)
	ServiceConfig := service_config_loader.ServiceConfigs[in.GetDataId()]
	dataId := nacosConn.GetDataId()
	_, ok := nacos_config_listener.NacosListedMap[dataId]
	if !ok {
		err := nacosConn.ServiceConfigListen()
		if err != nil {
			logs.Error(err)
			panic(err)
		} else {
			nacos_config_listener.NacosListedMap[dataId] = true
		}
	}

	response_, err := r.dubboHystrixServer("dubboServer", in, ServiceConfig)
	if err != nil {
		response.SetMessage(fmt.Sprintf("%s", err))
		panic(err)
	} else {
		response = response_
	}

	respCh <- response
}

func getNacosConn(in *io.RecRequest) nacos_config_listener.NacosConnConfig {
	//nacos listen need follow parms.
	nacosConn := nacos_config_listener.NacosConnConfig{}
	dataId := in.GetDataId()
	groupId := in.GetGroupId()
	namespaceId := in.GetNamespaceId()

	nacosConn.SetDataId(dataId)
	nacosConn.SetGroupId(groupId)
	nacosConn.SetNamespaceId(namespaceId)
	nacosConn.SetIp(ipAddr_)
	nacosConn.SetPort(port_)

	return nacosConn
}

func (r *DubbogoInferService) dubboHystrixServer(serverName string, in *io.RecRequest, ServiceConfig *service_config_loader.ServiceConfig) (*io.RecResponse, error) {
	response := &io.RecResponse{}
	response.SetCode(404)

	hystrixErr := hystrix.Do(serverName, func() error {
		// request recall / rank func.
		response_, err_ := r.RecommenderInfer(in, ServiceConfig)
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

		//INFO:降级和非降级用一个函数，会造成读写冲突
		itemList := in.GetItemList()
		in.SetRecallNum(int32(lowerRecallNum))
		in.SetItemList(itemList[:lowerRankNum])
		response_, err_ := r.RecommenderInferReduce(in, ServiceConfig)
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

func (r *DubbogoInferService) RecommenderInfer(in *io.RecRequest, ServiceConfig *service_config_loader.ServiceConfig) (*io.RecResponse, error) {
	response := &io.RecResponse{}
	response.SetCode(404)

	//build model by model_factory
	modelName := in.GetModelType()
	if modelName != "" {
		modelName = strings.ToLower(modelName)
	}

	modelfactory := model.ModelFactory{}
	request := getRequestParams(in)
	modelinfer, err := modelfactory.CreateInferModel(modelName, &request, ServiceConfig)
	if err != nil {
		return response, err
	}
	//close go2sky in dubbo service .
	//TODO: get go2sky config from config file.
	result, err := modelinfer.ModelInferNoSkywalking(nil) // d.dssmm.ModelInferNoSkywalking(nil)
	if err != nil {
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

func getRequestParams(in *io.RecRequest) io.RecRequest {
	request := io.RecRequest{}
	dataId := in.GetDataId()
	groupId := in.GetGroupId()
	namespaceId := in.GetNamespaceId()
	userId := in.GetUserId()
	itemList := in.GetItemList()

	request.SetDataId(dataId)
	request.SetGroupId(groupId)
	request.SetNamespaceId(namespaceId)
	request.SetUserId(userId)
	request.SetItemList(itemList)

	return request
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

func (r *DubbogoInferService) RecommenderInferReduce(in *io.RecRequest, ServiceConfig *service_config_loader.ServiceConfig) (*io.RecResponse, error) {
	response := &io.RecResponse{}
	response.SetCode(404)

	//build model by model_factory
	// modelName := in.GetModelType()
	// if modelName != "" {
	// 	modelName = strings.ToLower(modelName)
	// }

	modelName := "fm"

	modelfactory := model.ModelFactory{}
	request := getRequestParams(in)
	modelinfer, err := modelfactory.CreateInferModel(modelName, &request, ServiceConfig)
	if err != nil {
		return response, err
	}
	//close go2sky in dubbo service .
	//TODO: get go2sky config from config file.
	result, err := modelinfer.ModelInferNoSkywalking(nil) // d.dssmm.ModelInferNoSkywalking(nil)
	if err != nil {
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
