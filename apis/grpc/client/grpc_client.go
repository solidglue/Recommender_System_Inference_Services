package client

import (
	"fmt"
	grpc_server "infer-microservices/apis/grpc/server/api_gogofaster"
	"infer-microservices/utils/logs"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	// gRPC服务地址. 改成配置
	Address = "10.194.140.50:8221"
)

func req_grpc() {

	conn, err := grpc.Dial(Address, grpc.WithInsecure())
	if err != nil {
		logs.Error("GRPC dial failed")
	}
	defer conn.Close()

	client := grpc_server.NewGrpcRecommendServerServiceClient(conn)

	itemid_list := []string{"727407331", "745256420", "628549489", "637445275", "11111111111"} //"727407331","745256420","628549489","637445275","11111111111"

	req := &grpc_server.RecommendRequest{
		UserId: "13438935173",
		//ItemList: &grpc_server.StringList{itemid_list},
		ItemList: &grpc_server.StringList{Value: itemid_list},
	}

	fmt.Println("Address:", Address)

	fmt.Println("REQUEST:", req)

	//res, err := client.UinonRank(context.Background(), req)
	res, err := client.GrpcRecommendServer(context.Background(), req)

	if err != nil {
		logs.Error("GRPC req failed")
	}

	fmt.Println("RESPONSE:", res)

}
