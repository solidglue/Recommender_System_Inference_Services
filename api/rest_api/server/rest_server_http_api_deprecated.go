package server

import (
	"fmt"
	"infer-microservices/internal"
	"infer-microservices/internal/jwt"
	"infer-microservices/internal/logs"
	"infer-microservices/pkg/infer_services/rest_service"
	"net/http"
	"runtime"

	httpPlugin "github.com/SkyAPM/go2sky/plugins/http"
)

//TODO: test gin rest api, test dubbo restful api

type CallbackFunc func(w http.ResponseWriter, r *http.Request)

type HttpServiceApi struct {
	serverPort  uint
	maxCpuNum   int
	httpService *rest_service.HttpService
}

// set func
func (s *HttpServiceApi) SetServicePort(serverPort uint) {
	s.serverPort = serverPort
}

func (s *HttpServiceApi) SetMaxCpuNum(maxCpuNum int) {
	s.maxCpuNum = maxCpuNum
}

func (s *HttpServiceApi) SetRestService(httpService *rest_service.HttpService) {
	s.httpService = httpService
}

func (s *HttpServiceApi) restNoskywalkingServerRunner(path []string, CallbackFunc []CallbackFunc) error {
	for idx, p := range path {
		//http.HandleFunc(p, CallbackFunc [idx])
		http.Handle(p, jwt.JwtAuthMiddleware(http.HandlerFunc(CallbackFunc[idx])))
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

func (s *HttpServiceApi) restSkywalkingServerRunner(go2skyAddr string, serverName string, path []string, CallbackFunc []CallbackFunc) error {
	internal.SkywalkingTracer(go2skyAddr, serverName)

	sm, err := httpPlugin.NewServerMiddleware(internal.GetTracer())
	if err != nil {
		logs.Error("create server middleware error %v \n", err)
	}
	logs.Info("path:", path)
	logs.Info("CallbackFunc :", CallbackFunc)

	route := http.NewServeMux()
	for idx, p := range path {
		logs.Info("p CallbackFunc []:", p, CallbackFunc[idx])
		//route.HandleFunc(p, CallbackFunc [idx])
		route.Handle(p, jwt.JwtAuthMiddleware(http.HandlerFunc(CallbackFunc[idx])))

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
func (s *HttpServiceApi) ServiceStart() {
	paths := []string{
		"/login", "/infer",
	}

	CallbackFuncs := []CallbackFunc{
		jwt.AuthHandler, s.httpService.SyncRecommenderInfer,
	}

	if s.httpService.GetBaseService().GetSkywalkingWeatherOpen() {
		go2skyAddr := s.httpService.GetBaseService().GetSkywalkingIp() + ":" + fmt.Sprintf(":%d", s.httpService.GetBaseService().GetSkywalkingPort())
		go s.restSkywalkingServerRunner(go2skyAddr, s.httpService.GetBaseService().GetSkywalkingServerName(), paths, CallbackFuncs)
	} else {
		go s.restNoskywalkingServerRunner(paths, CallbackFuncs)
	}

	select {}
}
