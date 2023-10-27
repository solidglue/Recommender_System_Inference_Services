package server

import (
	"fmt"
	"infer-microservices/common"
	"infer-microservices/common/flags"
	"infer-microservices/utils/logs"
	"net/http"
	"runtime"

	httpPlugin "github.com/SkyAPM/go2sky/plugins/http"
)

//TODO: test gin rest api, test dubbo restful api

var restListenPort uint
var maxCpuNum int
var skywalkingWeatherOpen bool
var skywalkingIP string
var skywalkingPort int
var skywalkingServerName string
var NacosIP string
var NacosPort uint64

type WorkFunc func(w http.ResponseWriter, r *http.Request)
type HttpServer struct {
	ServerIP string
	Port     uint
}

func init() {
	flagFactory := flags.FlagFactory{}
	flagServiceConfig := flagFactory.FlagServiceConfigFactory()
	flagSkywalking := flagFactory.FlagSkywalkingFactory()

	restListenPort = *flagServiceConfig.GetServiceRestPort()
	maxCpuNum = *flagServiceConfig.GetServiceMaxCpuNum()
	skywalkingWeatherOpen = *flagSkywalking.GetSkywalkingWhetheropen()
	skywalkingIP = *flagSkywalking.GetSkywalkingIp()
	skywalkingPort = *flagSkywalking.GetSkywalkingPort()
	skywalkingServerName = *flagSkywalking.GetSkywalkingServername()
}

func NewHttpServer() *HttpServer {
	return &HttpServer{}
}

func (httpsvr *HttpServer) restNoskywalkingServerRunner(path []string, workFunc []WorkFunc) error {
	for idx, p := range path {
		http.HandleFunc(p, workFunc[idx])
	}

	cpuNum := runtime.NumCPU()
	if maxCpuNum <= cpuNum {
		cpuNum = maxCpuNum
	}

	runtime.GOMAXPROCS(cpuNum)
	logs.Info("cup num:", cpuNum)

	addr := fmt.Sprintf(":%d", restListenPort)
	err := http.ListenAndServe(addr, nil)
	if err == nil {
		logs.Error("server start succ ip:port ", addr)
		return err
	}

	return nil
}

func (httpsvr *HttpServer) restSkywalkingServerRunner(go2skyAddr string, serverName string, path []string, workFunc []WorkFunc) error {
	common.SkywalkingTracer(go2skyAddr, serverName)

	sm, err := httpPlugin.NewServerMiddleware(common.Tracer)
	if err != nil {
		logs.Error("create server middleware error %v \n", err)
	}
	logs.Info("path:", path)
	logs.Info("workFunc:", workFunc)

	route := http.NewServeMux()
	for idx, p := range path {
		logs.Info("p workFunc[]:", p, workFunc[idx])
		route.HandleFunc(p, workFunc[idx])
	}

	cpuNum := runtime.NumCPU()
	if maxCpuNum <= cpuNum {
		cpuNum = maxCpuNum
	}

	runtime.GOMAXPROCS(cpuNum)
	logs.Info("cup num:", cpuNum)

	addr := fmt.Sprintf(":%d", restListenPort)
	err = http.ListenAndServe(addr, sm(route))
	if err == nil {
		logs.Error("server start succ ip:port ", err)
		return nil
	}

	return nil
}

func RestServerRunner() {
	paths := []string{
		"/recall", "/rank",
	}

	restServer := &RestInferService{}
	workFunHandlers := []WorkFunc{
		restServer.restInferServer, restServer.restInferServer,
	}

	httpServer := NewHttpServer()
	if skywalkingWeatherOpen {
		go2skyAddr := skywalkingIP + ":" + fmt.Sprintf(":%d", skywalkingPort)
		go httpServer.restSkywalkingServerRunner(go2skyAddr, skywalkingServerName, paths, workFunHandlers)
	} else {
		go httpServer.restNoskywalkingServerRunner(paths, workFunHandlers)
	}

	select {}
}
