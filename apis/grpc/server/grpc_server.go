package server

//package main

import (
	"errors"
	"fmt"
	"infer-microservices/common/flags"
	"infer-microservices/cores/nacos_config"
	"infer-microservices/cores/service_config"
	"time"

	"infer-microservices/apis"
	grpc_api "infer-microservices/apis/grpc/server/api_gogofaster"
	"infer-microservices/apis/input_format"

	"infer-microservices/utils/logs"
	"net"
	"runtime"

	"github.com/afex/hystrix-go/hystrix"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	//"net/http"
)

//TODO: 传来的参数不固定，且枚举太多，考虑反射
//TODO: 补充grpc召回. will be remove
//TODO:comprass data
//INFO:recommend-go.proto

//TODO:1.启动监听nacos配置。2.服务注册到nacos
//"deepmodel_server/mg_online_predict/project/embedding_server"

var grpcListenPort uint
var maxCpuNum int
var skywalkingWeatherOpen bool
var lowerRecallNum int
var lowerRankNum int

var ipAddr_ string
var port_ uint64

var grpcInfer grpcInferInterface

// server is used to implement customer.CustomerServer.
type grpcRecommender struct {
}

func init() {
	// grpcListenPort = *flags.Grpc_server_port
	// maxCpuNum = *flags.Max_cpu_num
	// skywalkingWeatherOpen = *flags.Skywalking_whetheropen
	// lowerRecallNum = *flags.Hystrix_lowerRecallNum
	// lowerRankNum = *flags.Hystrix_lowerRankNum

	flagFactory := flags.FlagFactory{}
	flagServiceConfig := flagFactory.FlagServiceConfigFactory()
	flagSkywalking := flagFactory.FlagSkywalkingFactory()
	flagHystrix := flagFactory.FlagHystrixFactory()

	grpcListenPort = *flagServiceConfig.GetServiceGrpcPort()
	maxCpuNum = *flagServiceConfig.GetServiceMaxCpuNum()
	skywalkingWeatherOpen = *flagSkywalking.GetSkywalkingWhetheropen()
	lowerRecallNum = *flagHystrix.GetHystrixLowerRecallNum()
	lowerRankNum = *flagHystrix.GetHystrixLowerRankNum()

}

// TODO: rename to grpcRecommenderServer
func (g *grpcRecommender) GrpcRecommendServer(ctx context.Context, in *grpc_api.RecommendRequest) (*grpc_api.RecommendResponse, error) {

	//var err error

	resp_info := &grpc_api.RecommendResponse{
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
				//fmt.Println("context timeout exceeded")
				return resp_info, ctx.Err()
			case context.Canceled:
				//fmt.Println("context cancelled by force")
				return resp_info, ctx.Err()
			}
		case responseCh := <-respCh:
			resp_info = responseCh
			return resp_info, nil
		}
	}

}

func (g *grpcRecommender) grpcRecommenderServerContext(ctx context.Context, in *grpc_api.RecommendRequest, respCh chan *grpc_api.RecommendResponse) {

	defer func() {
		if info := recover(); info != nil {
			fmt.Println("触发了宕机", info)

		} else {
			//fmt.Println("程序正常退出")
		}
	}()

	resp_info := &grpc_api.RecommendResponse{
		Code: 404,
	}

	ServiceConfig := apis.ServiceConfigs[in.DataId]

	dataId := in.DataId
	groupId := in.GroupId
	namespaceId := in.Namespace

	request := apis.RecRequest{}
	request.SetDataId(dataId)
	request.SetGroupId(groupId)
	request.SetNamespaceId(namespaceId)
	request.SetUserId(in.UserId)
	request.SetRecallNum(in.RecallNum)
	request.SetItemList(in.ItemList.Value)

	ServiceConfig = apis.ServiceConfigs[dataId]

	nacosConn := nacos_config.NacosConnConfig{}
	nacosConn.SetDataId(dataId)
	nacosConn.SetGroupId(groupId)
	nacosConn.SetNamespaceId(namespaceId)
	nacosConn.SetIp(ipAddr_)
	nacosConn.SetPort(port_)

	_, ok := apis.NacosListedMap[dataId]
	if !ok {
		//注意：需求请求一次才会启动监听，启动空服务是不会监听的
		err := nacosConn.ServiceConfigListen()
		if err != nil {
			return
		} else {
			apis.NacosListedMap[dataId] = true
		}
	}

	//请求服务
	response_, err := g.grpcHystrixServer("grpcServer", &request, ServiceConfig)
	if err != nil {
		resp_info.Message = fmt.Sprintf("%s", err)
		return
	} else {
		resp_info = response_
	}

	respCh <- resp_info
}

func (r *grpcRecommender) grpcHystrixServer(serverName string, in *apis.RecRequest, ServiceConfig *service_config.ServiceConfig) (*grpc_api.RecommendResponse, error) {

	resp_info := &grpc_api.RecommendResponse{
		Code: 404,
	}

	hystrix.Do(serverName, func() error {

		// talk to other services
		response_, err := r.grpcRecommender(in, ServiceConfig)

		if err != nil {
			logs.Error(err)
		} else {
			resp_info = response_
		}

		return nil

	}, func(err error) error {

		// do this when services are down
		if err != nil {
			logs.Error(err)
		}

		itemList := in.GetItemList()
		in.SetRecallNum(int32(lowerRecallNum))
		in.SetItemList(itemList[:lowerRankNum])
		response_, err := r.grpcRecommender(in, ServiceConfig)

		if err != nil {
			logs.Error(err)
		} else {
			resp_info = response_
		}

		return nil
	})

	return resp_info, nil

}

func (g *grpcRecommender) grpcRecommender(in *apis.RecRequest, ServiceConfig *service_config.ServiceConfig) (*grpc_api.RecommendResponse, error) {

	resp_info := &grpc_api.RecommendResponse{
		Code: 404,
	}

	request := apis.RecRequest{}

	//TODO:	按照dubbo格式转换

	//modelType = MLMODEL, 就走MLReacll         //合并，用反射

	recaller := input_format.RecallInputFormat{}
	ranker := input_format.RankInputFormat{}

	modelType := in.GetModelType()
	if modelType == "recall" {

		dssm, err := recaller.InputCheckAndFormat(&request, ServiceConfig)
		if err != nil {
			logs.Error(err)
			return nil, err
		}
		grpcInfer = &recallServer{dssm}
		resp_info, err = grpcInfer.grpcInferServer()
		if err != nil {
			logs.Error(err)
			return nil, err
		}

	} else if modelType == "rank" {
		deepfm, err := ranker.InputCheckAndFormat(&request, ServiceConfig)
		if err != nil {
			logs.Error(err)
			return nil, err
		}

		grpcInfer = &rankServer{deepfm}
		resp_info, err = grpcInfer.grpcInferServer()
		if err != nil {
			logs.Error(err)
			return nil, err
		}

	} else {
		return resp_info, errors.New("wrong Strategy")
	}

	return resp_info, nil
}

func GrpcServerRunner(nacosIp string, nacosPort uint64) error {

	ipAddr_ = nacosIp
	port_ = nacosPort

	cpuNum := runtime.NumCPU()
	if maxCpuNum <= cpuNum {
		cpuNum = maxCpuNum
	}

	runtime.GOMAXPROCS(cpuNum)
	logs.Debug("cup num:", cpuNum)

	addr := fmt.Sprintf(":%d", grpcListenPort)

	lis, err := net.Listen("tcp", addr)
	//lis, err := net.Listen("tcp", ":8652")
	if err != nil {
		logs.Fatal("failed to listen: %v", err)
	} else {
		logs.Debug("listen to port:", addr)

	}

	s := grpc.NewServer()
	grpc_api.RegisterGrpcRecommendServerServiceServer(s, nil)

	s.Serve(lis)
	logs.Debug(" starting server ...", addr)

	if err != nil {
		logs.Error(err)
		return err
	}

	return nil

}
