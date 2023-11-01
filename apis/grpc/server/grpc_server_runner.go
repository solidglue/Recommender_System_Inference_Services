package server

import (
	"fmt"
	grpc_api "infer-microservices/apis/grpc/server/api_gogofaster"
	"infer-microservices/utils/logs"
	"net"
	"runtime"

	"google.golang.org/grpc"
)

type GrpcServer struct {
	serverPort            uint
	nacosIp               string
	nacosPort             uint
	maxCpuNum             int
	skywalkingWeatherOpen bool
	lowerRankNum          int
	lowerRecallNum        int
}

// set func

func (s *GrpcServer) SetServerPort(serverPort uint) {
	s.serverPort = serverPort
}

func (s *GrpcServer) SetNacosIp(nacosIp string) {
	s.nacosIp = nacosIp
}

func (s *GrpcServer) SetNacosPort(nacosPort uint) {
	s.nacosPort = nacosPort
}

func (s *GrpcServer) SetMaxCpuNum(maxCpuNum int) {
	s.maxCpuNum = maxCpuNum
}

func (s *GrpcServer) SetSkywalkingWeatherOpen(skywalkingWeatherOpen bool) {
	s.skywalkingWeatherOpen = skywalkingWeatherOpen
}

func (s *GrpcServer) SetLowerRankNum(lowerRankNum int) {
	s.lowerRankNum = lowerRankNum
}

func (s *GrpcServer) SetLowerRecallNum(lowerRecallNum int) {
	s.lowerRecallNum = lowerRecallNum
}

// @implement start infertace
func (s *GrpcServer) ServerStart() {
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
	grpc_api.RegisterGrpcRecommendServerServiceServer(gserver, s)
	gserver.Serve(lis)
	if err != nil {
		logs.Error(err)
	}
}
