package server

import (
	"context"
	"errors"
	"fmt"
	"infer-microservices/apis"
	"infer-microservices/apis/input_format"
	"infer-microservices/common/flags"
	"infer-microservices/cores/nacos_config"
	"infer-microservices/cores/service_config"
	"infer-microservices/utils/logs"
	"time"

	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports" // dubbogo 框架依赖，所有dubbogo进程都需要隐式引入一次
	"github.com/afex/hystrix-go/hystrix"
	hessian "github.com/apache/dubbo-go-hessian2"
)

//TODO: 传来的参数不固定，且枚举太多，考虑反射(性能差，慎用) https://blog.csdn.net/DkSakura/article/details/116588382
//TODO:反射应用点 - 如果api请求不规范，纠正一下？例如int类型传成了string类型，float类型与float64，
//，而不是直接拒绝，尤其是http请求时，用户手动输入不规范场景

//TODO:此处可以改用策略模式  https://zhuanlan.zhihu.com/p/392843448

//TODO:超时的话取消go协程执行，释放资源（虽然返回了，资源可能继续被占用。可能自带取消功能，需要确认

var ipAddr_ string
var port_ uint64
var lowerRecallNum int
var lowerRankNum int

var dubboInfer dubboInferInterface

type DubbogoInferService struct {
}

func init() {
	hessian.RegisterPOJO(&apis.RecRequest{})          // 注册传输结构到 hessian 库
	hessian.RegisterPOJO(&apis.RecResponse{})         // 注册传输结构到 hessian 库
	config.SetProviderService(&DubbogoInferService{}) // 注册服务提供者类，类名与配置文件中的 service 对应

	flagFactory := flags.FlagFactory{}
	flagHystrix := flagFactory.FlagHystrixFactory()
	lowerRecallNum = *flagHystrix.GetHystrixLowerRecallNum()
	lowerRankNum = *flagHystrix.GetHystrixLowerRankNum()

}

//INFO:dont remove
// // MethodMapper 定义方法名映射，从 Go 的方法名映射到 Java 小写方法名，只有 dubbo 协议服务接口才需要使用
// // go -> go 互通无需使用
// func (s *InferDubbogoService) MethodMapper() map[string]string {
// 	return map[string]string{
// 		"DubboRecommendServer": "dubboRecommendServer",
// 	}
// }

// 实现接口方法
func (r *DubbogoInferService) DubboRecommendServer(ctx context.Context, in *apis.RecRequest) (*apis.RecResponse, error) {

	response := &apis.RecResponse{}
	response.SetCode(404)

	//TODO：设置整体超时时间。用context熔断，hystrix降级
	ctx, cancelFunc := context.WithTimeout(ctx, time.Millisecond*150)
	defer cancelFunc()

	respCh := make(chan *apis.RecResponse, 100)

	go r.dubboRecommenderServerContext(ctx, in, respCh)

	for {
		select {
		case <-ctx.Done():
			switch ctx.Err() {
			case context.DeadlineExceeded:
				//fmt.Println("context timeout exceeded")
				return response, ctx.Err()
			case context.Canceled:
				//fmt.Println("context cancelled by force")
				return response, ctx.Err()
			}
		case responseCh := <-respCh:
			response = responseCh
			return response, nil
		}
	}

}

func (r *DubbogoInferService) dubboRecommenderServerContext(ctx context.Context, in *apis.RecRequest, respCh chan *apis.RecResponse) {

	response := &apis.RecResponse{}
	response.SetCode(404)

	dataId := in.GetDataId()
	groupId := in.GetGroupId()
	namespaceId := in.GetNamespaceId()
	ServiceConfig := apis.ServiceConfigs[in.GetDataId()]

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
	response_, err := r.dubboHystrixServer("dubboServer", in, ServiceConfig)
	if err != nil {
		response.SetMessage(fmt.Sprintf("%s", err))
		return
	} else {
		response = response_
	}

	respCh <- response

	//return response, nil //为避免dubbo框架处理错误，可不返回err，返回个nil就行  。20230220
}

func (r *DubbogoInferService) dubboHystrixServer(serverName string, in *apis.RecRequest, ServiceConfig *service_config.ServiceConfig) (*apis.RecResponse, error) {

	defer func() {
		if info := recover(); info != nil {
			fmt.Println("触发了宕机", info)

		} else {
			//fmt.Println("程序正常退出")
		}
	}()

	response := &apis.RecResponse{}
	response.SetCode(404)

	hystrix.Do(serverName, func() error {

		// talk to other services
		response_, err := r.dubboRecommender(in, ServiceConfig)

		if err != nil {
			logs.Error(err)
		} else {
			response = response_
		}

		return nil
	}, func(err error) error {

		//TODO: 此处降级如果超时，依然会卡住，需要结合context使用。或者直接熔断

		// do this when services are down
		if err != nil {
			logs.Error(err)
		}

		itemList := in.GetItemList()
		in.SetRecallNum(int32(lowerRecallNum))
		in.SetItemList(itemList[:lowerRankNum])
		response_, err := r.dubboRecommender(in, ServiceConfig)

		if err != nil {
			logs.Error(err)
		} else {
			response = response_
		}

		return nil
	})

	return response, nil

}

func (r *DubbogoInferService) dubboRecommender(in *apis.RecRequest, ServiceConfig *service_config.ServiceConfig) (*apis.RecResponse, error) {

	response := &apis.RecResponse{}
	response.SetCode(404)

	dataId := in.GetDataId()
	groupId := in.GetGroupId()
	namespaceId := in.GetNamespaceId()
	userId := in.GetUserId()
	itemList := in.GetItemList()

	request := apis.RecRequest{}
	request.SetDataId(dataId)
	request.SetGroupId(groupId)
	request.SetNamespaceId(namespaceId)
	request.SetUserId(userId)
	request.SetItemList(itemList)

	//TODO:此处可以改用策略模式  https://zhuanlan.zhihu.com/p/392843448

	modelType := in.GetModelType()

	if modelType == "recall" {

		recaller := input_format.RecallInputFormat{}
		dssm, err := recaller.InputCheckAndFormat(&request, ServiceConfig) //recaller.InputCheck(&request, ServiceConfig)
		if err != nil {
			logs.Error(err)
			return response, err
		}
		dubboInfer = &recallServer{dssm}
		response_, err := dubboInfer.dubboInferServer() //recaller.dubboRecallInferServer(dssm)
		if err != nil {
			logs.Error(err)
			return response, err
		} else {
			response = response_
		}

	} else if modelType == "rank" {
		ranker := input_format.RankInputFormat{}
		deepfm, err := ranker.InputCheckAndFormat(&request, ServiceConfig) //ranker.InputCheck(&request, ServiceConfig)
		if err != nil {
			logs.Error(err)
			return response, err
		}

		dubboInfer = &rankServer{deepfm}
		response_, err := dubboInfer.dubboInferServer()
		if err != nil {
			logs.Error(err)
			return response, err
		} else {
			response = response_
		}

	} else {
		err := errors.New("wrong Strategy")
		return response, err
	}

	return response, nil
}

// export DUBBO_GO_CONFIG_PATH=dubbogo.yml 运行前需要设置环境变量，指定配置文件位置
// 特别注意：start_dubbogo是私有函数，因此外边的main调用不到，一直报undefined错误，坑。。。。要改成Start_dubbogo
func DubboServerRunner(ipAddr string, port uint64, dubboConfFile string) {

	ipAddr_ = ipAddr
	port_ = port

	if err := config.Load(config.WithPath(dubboConfFile)); err != nil {
		panic(err)
	}
}
