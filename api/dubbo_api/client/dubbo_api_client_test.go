package client

import (
	"infer-microservices/pkg/services/io"
	"testing"
)

var dubboConfigFile string

func init() {
	dubboConfigFile = "../configs/dubbo/dubbogo_client_1.yml"
}

func TestDubboApiClient(t *testing.T) {
	itemList := []string{"1001", "1002", "1003", "1004"}

	req := io.RecRequest{}
	req.SetDataId("$dataid|$groupid") //nacos dataid|groupid
	req.SetUserId("$userid")          //userid
	req.SetItemList(itemList)         //rank items

	dubboApiClient := DubboApiClient{}
	dubboApiClient.setDubboConfigFile(dubboConfigFile)
	rsp, err := dubboApiClient.dubboServiceApiInfer(req)
	if err != nil {
		t.Fatal("service return err", err)
	}

	if len(rsp.GetData()) == 0 {
		t.Errorf("Expected recall_num or rank_num, but got 0")
	} else {
		t.Log("service rsp:", rsp)
	}
}
