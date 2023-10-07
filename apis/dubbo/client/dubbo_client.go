package client

import (
	"context"
	"fmt"
	"infer-microservices/apis"
	"infer-microservices/apis/dubbo/server/api"
	"infer-microservices/utils/logs"

	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
)

func requestDubboService() {
	itemList := []string{"7000000", "7000001", "7000002", "7000003"}

	req := apis.RecRequest{}
	req.SetDataId("dataid|groupid") //nacos dataid|groupid
	req.SetUserId("real userid")    //userid
	req.SetItemList(itemList)       //rank items

	// request dubbo infer service
	rsp, err := api.DubbogoInferServiceClient.DubboRecommendServer(context.TODO(), &req)
	if err != nil {
		logs.Error(err)
	}

	for i := 0; i < len(rsp.GetData()); i++ {
		fmt.Println(i, req.GetUserId(), rsp.GetData()[i])
	}
}

// TODO: change 2 unit test.
func main() {
	// export DUBBO_GO_CONFIG_PATH=dubbogo.yml or load it in code.
	if err := config.Load(config.WithPath("conf/dubbogo_client.yml")); err != nil {
		panic(err)
	}

	requestDubboService() //recall or rank depends on nacos config file.
}
