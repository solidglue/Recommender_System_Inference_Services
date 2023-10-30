package common

import (
	"infer-microservices/common/flags"

	"github.com/afex/hystrix-go/hystrix"
)

var hystrixTimeout int
var hystrixMaxConcurrentRequests int
var hystrixRequestVolumeThreshold int
var hystrixSleepWindow int
var hystrixErrorPercentThreshold int

func init() {
	flagFactory := flags.FlagFactory{}
	flagHystrix := flagFactory.CreateFlagHystrix()

	timeout := *flagHystrix.GetHystrixTimeoutMs()
	maxConcurrentRequests := *flagHystrix.GetHystrixMaxConcurrentRequests()
	requestVolumeThreshold := *flagHystrix.GetHystrixErrorPercentThreshold()
	sleepWindow := *flagHystrix.GetHystrixSleepWindow()
	errorPercentThreshold := *flagHystrix.GetHystrixErrorPercentThreshold()

	//all ConfigureCommands use the same conf, we can use different conf if we need.
	hystrix.ConfigureCommand("dubboServer", hystrix.CommandConfig{
		Timeout:                timeout,
		MaxConcurrentRequests:  maxConcurrentRequests,
		RequestVolumeThreshold: requestVolumeThreshold,
		SleepWindow:            sleepWindow,
		ErrorPercentThreshold:  errorPercentThreshold,
	})

	hystrix.ConfigureCommand("restServer", hystrix.CommandConfig{
		Timeout:                timeout,
		MaxConcurrentRequests:  maxConcurrentRequests,
		RequestVolumeThreshold: requestVolumeThreshold,
		SleepWindow:            sleepWindow,
		ErrorPercentThreshold:  errorPercentThreshold,
	})

	hystrix.ConfigureCommand("grpcServer", hystrix.CommandConfig{
		Timeout:                timeout,
		MaxConcurrentRequests:  maxConcurrentRequests,
		RequestVolumeThreshold: requestVolumeThreshold,
		SleepWindow:            sleepWindow,
		ErrorPercentThreshold:  errorPercentThreshold,
	})

}
