package client

import (
	"infer-microservices/pkg/logs"
	grpc_service "infer-microservices/pkg/services/grpc_service"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type GrpcApiClient struct {
	grpcAddress string
}

func (d *GrpcApiClient) setGrpcAddress(grpcAddress string) {
	d.grpcAddress = grpcAddress
}

func (g *GrpcApiClient) grpcServiceApiInfer(req *grpc_service.RecommendRequest) (*grpc_service.RecommendResponse, error) {
	conn, err := grpc.Dial(g.grpcAddress, grpc.WithInsecure())
	if err != nil {
		logs.Error("GRPC dial failed")
	}
	defer conn.Close()

	client := grpc_service.NewRecommenderInferServiceClient(conn)
	rsp, err := client.RecommenderInfer(context.Background(), req)
	if err != nil {
		logs.Error(err)
		return nil, err
	}

	return rsp, nil
}
