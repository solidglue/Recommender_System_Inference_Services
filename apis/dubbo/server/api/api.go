package api

import (
	"context"
	"infer-microservices/apis"

	"dubbo.apache.org/dubbo-go/v3/config"
	hessian "github.com/apache/dubbo-go-hessian2"
	//"time"
)

// https://www.w3cschool.cn/dubbo/languages-golang-dubbo-go-30-quickstart-quickstart_dubbo.html

func init() {
	hessian.RegisterPOJO(&apis.RecRequest{}) // 注册传输结构到 hessian 库
	//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	//没注册这个接口，卡了3天。。。。。官方样例返回和输入都是同一个
	hessian.RegisterPOJO(&apis.RecResponse{}) // 注册传输结构到 hessian 库

	// 注册客户端存根类到框架，实例化客户端接口指针 userProvider
	config.SetConsumerService(DubbogoInferServiceClient)
}

var (
	DubbogoInferServiceClient = &DubbogoInferService{} // 客户端指针
)

// 2。 定义客户端存根类：UserProvider
// 特征注意，返回RecResponse，报错，说是timeout或tcp错误。改成string就正常了？？？？？？？？？
type DubbogoInferService struct {
	// dubbo标签，用于适配go侧客户端大写方法名 -> java侧小写方法名，只有 dubbo 协议客户端才需要使用
	DubboRecommendServer func(ctx context.Context, req *apis.RecRequest) (*apis.RecResponse, error) //`dubbo:"getUser"`

}
