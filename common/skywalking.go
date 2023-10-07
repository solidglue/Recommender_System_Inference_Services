package common

import (
	"infer-microservices/utils/logs"

	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/reporter"
)

// go2sky
var report go2sky.Reporter
var Tracer *go2sky.Tracer
var errReport error
var errTracer error

func SkywalkingTracer(go2skyAddr string, serverName string) {

	report, errReport = reporter.NewGRPCReporter(go2skyAddr)
	Tracer, errTracer = go2sky.NewTracer(serverName, go2sky.WithReporter(report)) //// tracer类型    go2sky.Tracer

	if errReport != nil {
		logs.Error("crate grpc reporter error: %v \n", errReport)
	}

	if errTracer != nil {
		logs.Error("crate tracer error: %v \n", errTracer)
	}

}
