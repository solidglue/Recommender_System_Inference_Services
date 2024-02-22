package client

import (
	grpc_service "infer-microservices/pkg/infer_services/grpc_service"
	"testing"
)

var grpcAddress string

func init() {
	grpcAddress = "10.10.10.10:8888"
}

func TestDubboApiClient(t *testing.T) {
	itemList := []string{"1001", "1002", "1003", "1004"}
	req := &grpc_service.RecommendRequest{
		UserId:   "$userid",
		ItemList: &grpc_service.StringList{Value: itemList},
	}

	grpcClient := GrpcApiClient{}
	grpcClient.setGrpcAddress(grpcAddress)
	rsp, err := grpcClient.grpcServiceApiInfer(req)
	if err != nil {
		t.Fatal("service return err", err)
	}

	if rsp.Data.Size() == 0 {
		t.Errorf("Expected recall_num or rank_num, but got 0")
	} else {
		t.Log(rsp.Data.Iteminfo_)
	}
}
