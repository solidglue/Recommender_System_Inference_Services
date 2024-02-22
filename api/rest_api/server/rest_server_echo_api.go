package server

import (
	"fmt"
	"infer-microservices/internal"
	"infer-microservices/internal/jwt"
	"infer-microservices/internal/logs"
	"infer-microservices/pkg/infer_services/rest_service"
	"runtime"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

//TODO: test gin rest api, test dubbo restful api

var echoApi *echo.Echo
var skywalkingOpen bool
var skywalkingAddr string
var skywalkingServerName string

type EchoServiceApi struct {
	serverPort  uint
	maxCpuNum   int
	echoService *rest_service.EchoService
}

func init() {
	echoApi = echo.New()
	echoApi.Debug = true
	echoApi.Use(middleware.Recover()) // 主要用于拦截panic错误并且在控制台打印错误日志，避免echo程序直接崩溃
	echoApi.Use(middleware.Logger())  // Logger中间件主要用于打印http请求日志
	echoApi.Use(middleware.RequestID())
}

// 全局使用echo对象
func GetEcho() *echo.Echo {
	return echoApi
}

// set func
func (s *EchoServiceApi) SetServicePort(serverPort uint) {
	s.serverPort = serverPort
}

func (s *EchoServiceApi) SetMaxCpuNum(maxCpuNum int) {
	s.maxCpuNum = maxCpuNum
}

func (s *EchoServiceApi) SetRestService(echoService *rest_service.EchoService) {
	s.echoService = echoService
}

// skywalkingMiddleware
func skywalkingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if skywalkingOpen {
			internal.SkywalkingTracer(skywalkingAddr, skywalkingServerName)
		}

		return next(c)
	}
}

// @implement start infertace
func (s *EchoServiceApi) ServiceStart() {
	echoApi.POST("/login2", jwt.Login)

	// Restricted group
	r := echoApi.Group("/infer2")

	// Configure middleware with the custom claims type
	config := middleware.JWTConfig{
		Claims:     &jwt.JwtCustomClaims{},
		SigningKey: []byte("secret"),
	}
	r.Use(middleware.JWTWithConfig(config))
	echoApi.POST("/infer2", s.echoService.SyncRecommenderInfer)

	skywalkingOpen = s.echoService.GetBaseService().GetSkywalkingWeatherOpen()
	skywalkingAddr = s.echoService.GetBaseService().GetSkywalkingIp() + ":" + fmt.Sprintf(":%d", s.echoService.GetBaseService().GetSkywalkingPort())
	skywalkingServerName = s.echoService.GetBaseService().GetSkywalkingServerName()
	echoApi.Use(skywalkingMiddleware)

	cpuNum := runtime.NumCPU()
	if s.maxCpuNum <= cpuNum {
		cpuNum = s.maxCpuNum
	}
	runtime.GOMAXPROCS(cpuNum)
	logs.Info("cup num:", cpuNum)

	addr := fmt.Sprintf(":%d", s.serverPort)
	echoApi.Logger.Fatal(echoApi.Start(addr))
	logs.Info("server addr:", addr)
}
