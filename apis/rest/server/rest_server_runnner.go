package server

import (
	"fmt"
	"infer-microservices/common"
	"infer-microservices/utils/logs"
	"net/http"
	"runtime"

	httpPlugin "github.com/SkyAPM/go2sky/plugins/http"
)

//TODO: test gin rest api, test dubbo restful api

type InferFunc func(w http.ResponseWriter, r *http.Request)

type HttpServer struct {
	serverPort            uint
	nacosIp               string
	nacosPort             uint
	maxCpuNum             int
	skywalkingWeatherOpen bool
	skywalkingIp          string
	skywalkingPort        uint
	skywalkingServerName  string
	lowerRankNum          int
	lowerRecallNum        int
}

// set func

func (s *HttpServer) SetServerPort(serverPort uint) {
	s.serverPort = serverPort
}

func (s *HttpServer) SetNacosIp(nacosIp string) {
	s.nacosIp = nacosIp
}

func (s *HttpServer) SetNacosPort(nacosPort uint) {
	s.nacosPort = nacosPort
}

func (s *HttpServer) SetMaxCpuNum(maxCpuNum int) {
	s.maxCpuNum = maxCpuNum
}

func (s *HttpServer) SetSkywalkingWeatherOpen(skywalkingWeatherOpen bool) {
	s.skywalkingWeatherOpen = skywalkingWeatherOpen
}

func (s *HttpServer) SetSkywalkingIp(skywalkingIp string) {
	s.skywalkingIp = skywalkingIp
}

func (s *HttpServer) SetSkywalkingPort(skywalkingPort uint) {
	s.skywalkingPort = skywalkingPort
}

func (s *HttpServer) SetSkywalkingServerName(skywalkingServerName string) {
	s.skywalkingServerName = skywalkingServerName
}

func (s *HttpServer) SetLowerRankNum(lowerRankNum int) {
	s.lowerRankNum = lowerRankNum
}

func (s *HttpServer) SetLowerRecallNum(lowerRecallNum int) {
	s.lowerRecallNum = lowerRecallNum
}

func (s *HttpServer) restNoskywalkingServerRunner(path []string, InferFunc []InferFunc) error {
	for idx, p := range path {
		http.HandleFunc(p, InferFunc[idx])
	}

	cpuNum := runtime.NumCPU()
	if s.maxCpuNum <= cpuNum {
		cpuNum = s.maxCpuNum
	}

	runtime.GOMAXPROCS(cpuNum)
	logs.Info("cup num:", cpuNum)

	addr := fmt.Sprintf(":%d", s.serverPort)
	err := http.ListenAndServe(addr, nil)
	if err == nil {
		logs.Error("server start succ ip:port ", addr)
		return err
	}

	return nil
}

func (s *HttpServer) restSkywalkingServerRunner(go2skyAddr string, serverName string, path []string, InferFunc []InferFunc) error {
	common.SkywalkingTracer(go2skyAddr, serverName)

	sm, err := httpPlugin.NewServerMiddleware(common.Tracer)
	if err != nil {
		logs.Error("create server middleware error %v \n", err)
	}
	logs.Info("path:", path)
	logs.Info("InferFunc:", InferFunc)

	route := http.NewServeMux()
	for idx, p := range path {
		logs.Info("p InferFunc[]:", p, InferFunc[idx])
		route.HandleFunc(p, InferFunc[idx])
	}

	cpuNum := runtime.NumCPU()
	if s.maxCpuNum <= cpuNum {
		cpuNum = s.maxCpuNum
	}

	runtime.GOMAXPROCS(cpuNum)
	logs.Info("cup num:", cpuNum)

	addr := fmt.Sprintf(":%d", s.serverPort)
	err = http.ListenAndServe(addr, sm(route))
	if err == nil {
		logs.Error("server start succ ip:port ", err)
		return nil
	}

	return nil
}

// @implement start infertace
func (s *HttpServer) ServerStart() {
	paths := []string{
		"/infer",
	}

	InferFuncs := []InferFunc{
		s.restInferServer,
	}

	if s.skywalkingWeatherOpen {
		go2skyAddr := s.skywalkingIp + ":" + fmt.Sprintf(":%d", s.skywalkingPort)
		go s.restSkywalkingServerRunner(go2skyAddr, s.skywalkingServerName, paths, InferFuncs)
	} else {
		go s.restNoskywalkingServerRunner(paths, InferFuncs)
	}

	select {}
}
