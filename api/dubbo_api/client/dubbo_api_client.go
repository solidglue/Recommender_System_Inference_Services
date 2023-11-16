package client

import (
	"context"
	"fmt"
	dubbo_server "infer-microservices/api/dubbo_api/server"
	"infer-microservices/pkg/logs"
	"infer-microservices/pkg/services/io"

	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
)

//INFO:use to test serivce.

func requestDubboService() {
	itemList := []string{"80000001", "80000002", "80000003", "80000004"}

	req := io.RecRequest{}
	req.SetDataId("$dataid|$groupid") //nacos dataid|groupid
	req.SetUserId("$userid")          //userid
	req.SetItemList(itemList)         //rank items

	// request dubbo infer service. use to test serivce.
	rsp, err := dubbo_server.DubboServiceApiClient.RecommenderInfer(context.TODO(), &req)
	if err != nil {
		logs.Error(err)
	}

	for i := 0; i < len(rsp.GetData()); i++ {
		fmt.Println(i, req.GetUserId(), rsp.GetData()[i])
	}
}

// TODO: change to unit test.
func main() {
	// export DUBBO_GO_CONFIG_PATH=dubbogo.yml or load it in code.
	if err := config.Load(config.WithPath("conf/dubbogo_client.yml")); err != nil {
		panic(err)
	}

	requestDubboService() //recall or rank depends on nacos config file.
}
