package server

import (
	"fmt"
	grpc_api "infer-microservices/apis/grpc/server/api_gogofaster"
	"infer-microservices/utils/logs"
	"net"
	"runtime"

	"google.golang.org/grpc"
)

// runner
func GrpcServerRunner(nacosIp string, nacosPort uint64) error {
	ipAddr_ = nacosIp
	port_ = nacosPort
	cpuNum := runtime.NumCPU()
	if maxCpuNum <= cpuNum {
		cpuNum = maxCpuNum
	}
	runtime.GOMAXPROCS(cpuNum)
	logs.Info("cup num:", cpuNum)

	addr := fmt.Sprintf(":%d", grpcListenPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logs.Fatal("failed to listen: %v", err)
		panic(err)
	} else {
		logs.Info("listen to port:", addr)
	}

	s := grpc.NewServer()
	grpc_api.RegisterGrpcRecommendServerServiceServer(s, &GrpcInferService{})
	s.Serve(lis)
	if err != nil {
		logs.Error(err)
		return err
	}

	return nil
}
