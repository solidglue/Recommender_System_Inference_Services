package internal

import (
	"infer-microservices/internal/logs"

	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/reporter"
)

// go2sky
var report go2sky.Reporter
var tracer *go2sky.Tracer
var errReport error
var errTracer error

func GetTracer() *go2sky.Tracer {
	return tracer
}

func SkywalkingTracer(go2skyAddr string, serverName string) {
	report, errReport = reporter.NewGRPCReporter(go2skyAddr)
	tracer, errTracer = go2sky.NewTracer(serverName, go2sky.WithReporter(report))
	if errReport != nil {
		logs.Error("crate grpc reporter error: %v \n", errReport)
	}
	if errTracer != nil {
		logs.Error("crate tracer error: %v \n", errTracer)
	}
}
