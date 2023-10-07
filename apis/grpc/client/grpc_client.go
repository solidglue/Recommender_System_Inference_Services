package client

import (
	"fmt"
	grpc_server "infer-microservices/apis/grpc/server/api_gogofaster"
	"infer-microservices/utils/logs"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	// gRPC addrs.
	Address = "10.194.140.50:8221"
)

func requestGrpcService() {
	conn, err := grpc.Dial(Address, grpc.WithInsecure())
	if err != nil {
		logs.Error("GRPC dial failed")
	}
	defer conn.Close()

	itemList := []string{"7000000", "7000001", "7000002", "7000003"}
	req := &grpc_server.RecommendRequest{
		UserId: "real userid",
		//ItemList: &grpc_server.StringList{itemid_list},
		ItemList: &grpc_server.StringList{Value: itemList},
	}

	fmt.Println("REQUEST:", req)
	client := grpc_server.NewGrpcRecommendServerServiceClient(conn)
	res, err := client.GrpcRecommendServer(context.Background(), req)
	if err != nil {
		logs.Error("GRPC request failed")
	}
	fmt.Println("RESPONSE:", res)
}

// TODO: change 2 unit test.
func main() {
	requestGrpcService()
}
