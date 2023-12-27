package server

import (
	"fmt"
	"infer-microservices/internal/logs"
	grpc_service "infer-microservices/pkg/services/grpc_service"
	"net"
	"runtime"

	"google.golang.org/grpc"
)

type GrpcServiceApi struct {
	maxCpuNum   int
	serverPort  uint
	grpcService *grpc_service.GrpcService
}

func (s *GrpcServiceApi) SetMaxCpuNum(maxCpuNum int) {
	s.maxCpuNum = maxCpuNum
}

func (s *GrpcServiceApi) SetServicePort(serverPort uint) {
	s.serverPort = serverPort
}

func (s *GrpcServiceApi) SetGrpcService(grpcService *grpc_service.GrpcService) {
	s.grpcService = grpcService
}

// @implement start infertace
func (s *GrpcServiceApi) ServiceStart() {
	cpuNum := runtime.NumCPU()
	if s.maxCpuNum <= cpuNum {
		cpuNum = s.maxCpuNum
	}
	runtime.GOMAXPROCS(cpuNum)
	logs.Info("cup num:", cpuNum)

	addr := fmt.Sprintf(":%d", s.serverPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logs.Fatal("failed to listen: %v", err)
		panic(err)
	} else {
		logs.Info("listen to port:", addr)
	}

	gserver := grpc.NewServer()
	grpc_service.RegisterRecommenderInferServiceServer(gserver, s.grpcService)
	gserver.Serve(lis)
	if err != nil {
		logs.Error(err)
	}
}
