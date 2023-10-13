package client

import (
	"fmt"
	grpc_server "infer-microservices/apis/grpc/server/api_gogofaster"
	"infer-microservices/utils/logs"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

//INFO:use to test serivce.

const (
	// gRPC addrs.
	//TODO: //TODO: get  config from config file.
	Address = "10.10.10.10:8888"
)

func requestGrpcService() {
	conn, err := grpc.Dial(Address, grpc.WithInsecure())
	if err != nil {
		logs.Error("GRPC dial failed")
	}
	defer conn.Close()

	itemList := []string{"80000002", "80000002", "80000003", "80000004"}
	req := &grpc_server.RecommendRequest{
		UserId:   "$userid",
		ItemList: &grpc_server.StringList{Value: itemList},
	}

	fmt.Println("request:", req)
	client := grpc_server.NewGrpcRecommendServerServiceClient(conn)
	res, err := client.GrpcRecommendServer(context.Background(), req)
	if err != nil {
		logs.Error("GRPC request failed")
	}
	fmt.Println("response:", res)
}

// TODO: change to unit test.
func main() {
	requestGrpcService()
}
