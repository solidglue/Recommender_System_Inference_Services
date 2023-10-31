package server

import (
	"fmt"
	"infer-microservices/common/flags"
	"infer-microservices/cores/model"
	"infer-microservices/cores/nacos_config_listener"
	"infer-microservices/cores/service_config_loader"
	"strings"
	"sync"
	"time"

	grpc_api "infer-microservices/apis/grpc/server/api_gogofaster"
	"infer-microservices/apis/io"

	"infer-microservices/utils/logs"

	"github.com/afex/hystrix-go/hystrix"
	"golang.org/x/net/context"
)

var grpcListenPort uint
var maxCpuNum int
var skywalkingWeatherOpen bool
var lowerRecallNum int
var lowerRankNum int
var ipAddr_ string
var port_ uint64
var inferWg sync.WaitGroup

// server is used to implement customer.CustomerServer.
type GrpcInferService struct {
}

func init() {
	flagFactory := flags.FlagFactory{}
	flagServiceConfig := flagFactory.CreateFlagServiceConfig()
	flagSkywalking := flagFactory.CreateFlagSkywalking()
	flagHystrix := flagFactory.CreateFlagHystrix()

	grpcListenPort = *flagServiceConfig.GetServiceGrpcPort()
	maxCpuNum = *flagServiceConfig.GetServiceMaxCpuNum()
	skywalkingWeatherOpen = *flagSkywalking.GetSkywalkingWhetheropen()
	lowerRecallNum = *flagHystrix.GetHystrixLowerRecallNum()
	lowerRankNum = *flagHystrix.GetHystrixLowerRankNum()
}

// INFO: implement grpc func which defined by proto.
func (g *GrpcInferService) GrpcRecommendServer(ctx context.Context, in *grpc_api.RecommendRequest) (*grpc_api.RecommendResponse, error) {
	//INFO: set timeout by context, degraded service by hystix.
	response := &grpc_api.RecommendResponse{
		Code: 404,
	}

	ctx, cancelFunc := context.WithTimeout(ctx, time.Millisecond*150)
	defer cancelFunc()

	respCh := make(chan *grpc_api.RecommendResponse, 100)
	go g.grpcRecommenderServerContext(ctx, in, respCh)
	for {
		select {
		case <-ctx.Done():
			switch ctx.Err() {
			case context.DeadlineExceeded:
				return response, ctx.Err()
			case context.Canceled:
				return response, ctx.Err()
			}
		case responseCh := <-respCh:
			response = responseCh
			return response, nil
		}
	}
}

func getGrpcRequestParams(in *grpc_api.RecommendRequest) io.RecRequest {
	request := io.RecRequest{}
	request.SetDataId(in.GetDataId())
	request.SetGroupId(in.GetGroupId())
	request.SetNamespaceId(in.GetNamespace())
	request.SetUserId(in.UserId)
	request.SetRecallNum(in.RecallNum)
	request.SetItemList(in.ItemList.Value)

	return request
}

func getNacosConn(in *grpc_api.RecommendRequest) nacos_config_listener.NacosConnConfig {
	//nacos listen need follow parms.
	nacosConn := nacos_config_listener.NacosConnConfig{}
	nacosConn.SetDataId(in.GetDataId())
	nacosConn.SetGroupId(in.GetGroupId())
	nacosConn.SetNamespaceId(in.GetNamespace())
	nacosConn.SetIp(ipAddr_)
	nacosConn.SetPort(port_)

	return nacosConn
}

func (g *GrpcInferService) grpcRecommenderServerContext(ctx context.Context, in *grpc_api.RecommendRequest, respCh chan *grpc_api.RecommendResponse) {

	defer func() {
		if info := recover(); info != nil {
			fmt.Println("panic", info)
		} //else {
		//fmt.Println("")
		//}
	}()

	response := &grpc_api.RecommendResponse{
		Code: 404,
	}

	nacosConn := getNacosConn(in)
	dataId := in.GetDataId()
	ServiceConfig := service_config_loader.ServiceConfigs[dataId]
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
	request := getGrpcRequestParams(in)
	response_, err := g.grpcHystrixServer("grpcServer", &request, ServiceConfig)
	if err != nil {
		response.Message = fmt.Sprintf("%s", err)
		panic(err)
	} else {
		response = response_
	}
	respCh <- response
}

func (r *GrpcInferService) grpcHystrixServer(serverName string, in *io.RecRequest, ServiceConfig *service_config_loader.ServiceConfig) (*grpc_api.RecommendResponse, error) {
	response := &grpc_api.RecommendResponse{
		Code: 404,
	}

	hystrixErr := hystrix.Do(serverName, func() error {
		// request recall / rank func.
		response_, err := r.RecommenderInfer(in, ServiceConfig)
		if err != nil {
			logs.Error(err)
			return err
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

func (g *GrpcInferService) RecommenderInfer(in *io.RecRequest, ServiceConfig *service_config_loader.ServiceConfig) (*grpc_api.RecommendResponse, error) {
	response := &grpc_api.RecommendResponse{
		Code: 404,
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

	var responseTf map[string]interface{}
	if skywalkingWeatherOpen {
		responseTf, err = modelinfer.ModelInferSkywalking(nil)
	} else {
		responseTf, err = modelinfer.ModelInferNoSkywalking(nil)
	}

	if err != nil {
		logs.Error("request tfserving fail:", responseTf)
		return response, err
	}

	result := make([]*grpc_api.ItemInfo, 0)
	resultList := responseTf["data"].([]map[string]interface{})
	rankCh := make(chan *grpc_api.ItemInfo, len(resultList))
	for i := 0; i < len(resultList); i++ {
		inferWg.Add(1)
		go formatGrpcResponse(resultList[i], rankCh)
	}

	inferWg.Wait()
	close(rankCh)
	for itemScore := range rankCh {
		result = append(result, itemScore)
	}

	response = &grpc_api.RecommendResponse{
		Code:    200,
		Message: "success",
		Data: &grpc_api.ItemInfoList{
			Iteminfo_: result,
		},
	}

	return response, nil
}

func formatGrpcResponse(itemScore map[string]interface{}, rankCh chan *grpc_api.ItemInfo) {
	defer inferWg.Done()

	itemid := itemScore["itemid"].(string)
	score := float32(itemScore["score"].(float64))
	itemInfo := &grpc_api.ItemInfo{
		Itemid: itemid,
		Score:  score,
	}

	rankCh <- itemInfo
}

func (g *GrpcInferService) RecommenderInferReduce(in *io.RecRequest, ServiceConfig *service_config_loader.ServiceConfig) (*grpc_api.RecommendResponse, error) {
	response := &grpc_api.RecommendResponse{
		Code: 404,
	}

	//build model by model_factory
	// modelName := in.GetModelType()
	// if modelName != "" {
	// 	modelName = strings.ToLower(modelName)
	// }

	modelName := "fm"

	modelfactory := model.ModelFactory{}
	modelinfer, err := modelfactory.CreateInferModel(modelName, in, ServiceConfig)
	if err != nil {
		return response, err
	}

	var responseTf map[string]interface{}
	if skywalkingWeatherOpen {
		responseTf, err = modelinfer.ModelInferSkywalking(nil)
	} else {
		responseTf, err = modelinfer.ModelInferNoSkywalking(nil)
	}

	if err != nil {
		logs.Error("request tfserving fail:", responseTf)
		return response, err
	}

	result := make([]*grpc_api.ItemInfo, 0)
	resultList := responseTf["data"].([]map[string]interface{})
	rankCh := make(chan *grpc_api.ItemInfo, len(resultList))
	for i := 0; i < len(resultList); i++ {
		inferWg.Add(1)
		go formatGrpcResponse(resultList[i], rankCh)
	}

	inferWg.Wait()
	close(rankCh)
	for itemScore := range rankCh {
		result = append(result, itemScore)
	}

	response = &grpc_api.RecommendResponse{
		Code:    200,
		Message: "success",
		Data: &grpc_api.ItemInfoList{
			Iteminfo_: result,
		},
	}

	return response, nil
}
