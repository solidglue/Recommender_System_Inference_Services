package server

import grpc_api "infer-microservices/apis/grpc/server/api_gogofaster"

type grpcInferInterface interface {
	grpcInferServer() (*grpc_api.RecommendResponse, error)
}
