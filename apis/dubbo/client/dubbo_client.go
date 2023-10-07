package client

import (
	"context"
	"fmt"

	"infer-microservices/apis"
	"infer-microservices/apis/dubbo/server/api"

	_ "dubbo.apache.org/dubbo-go/v3/imports" // dubbogo 框架依赖，所有dubbogo进程都需要隐式引入一次
	//"deepmodel_server/mg_micro_server/dubbogo/api"
)

func req1() {

	item_list := []string{"677387947", "608548125", "608604182", "637445275", "11111111111"}
	// req := apis.RecRequest{

	// 	DataId:   "dssm-120|recengine-model4g",
	// 	UserId:   "00149a82cd13ce6ed0c73bc5f522b5a4",
	// 	ItemList: item_list,
	// }
	req := apis.RecRequest{}
	req.SetDataId("dssm-120|recengine-model4g")
	req.SetUserId("00149a82cd13ce6ed0c73bc5f522b5a4")
	req.SetItemList(item_list)

	// 发起调用
	rsp, err := api.DubbogoInferServiceClient.DubboRecommendServer(context.TODO(), &req)
	if err != nil {
		fmt.Println(">>>>>>>>>>ERRRRRRRRRRRROR>>>>>>>>>>>>", err)
	}
	//logger.Infof("response result: %+v", rsp)

	for i := 0; i < len(rsp.GetData()); i++ {
		fmt.Println(i, req.GetUserId(), rsp.GetData()[i])
	}

}

// // export DUBBO_GO_CONFIG_PATH=dubbogo.yml 运行前需要设置环境变量，指定配置文件位置
// func main() {
// 	// 启动框架
// 	//if err := config.Load(); err != nil{     //err := config.Load(config.WithPath("./conf/dubbo.yml"))
// 	if err := config.Load(config.WithPath("dubbogo.yml")); err != nil {
// 		panic(err)
// 	}

// 	req1()
// 	req2()

// }
