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

	hystrixTimeout = *flagHystrix.GetHystrixTimeoutMs()
	hystrixMaxConcurrentRequests = *flagHystrix.GetHystrixMaxConcurrentRequests()
	hystrixRequestVolumeThreshold = *flagHystrix.GetHystrixErrorPercentThreshold()
	hystrixSleepWindow = *flagHystrix.GetHystrixSleepWindow()
	hystrixErrorPercentThreshold = *flagHystrix.GetHystrixErrorPercentThreshold()

	//all ConfigureCommands use the same conf, we can use different conf if we need.
	hystrix.ConfigureCommand("dubboServer", hystrix.CommandConfig{
		Timeout:                hystrixTimeout,
		MaxConcurrentRequests:  hystrixMaxConcurrentRequests,
		RequestVolumeThreshold: hystrixRequestVolumeThreshold,
		SleepWindow:            hystrixSleepWindow,
		ErrorPercentThreshold:  hystrixErrorPercentThreshold,
	})

	hystrix.ConfigureCommand("restServer", hystrix.CommandConfig{
		Timeout:                hystrixTimeout,
		MaxConcurrentRequests:  hystrixMaxConcurrentRequests,
		RequestVolumeThreshold: hystrixRequestVolumeThreshold,
		SleepWindow:            hystrixSleepWindow,
		ErrorPercentThreshold:  hystrixErrorPercentThreshold,
	})

	hystrix.ConfigureCommand("grpcServer", hystrix.CommandConfig{
		Timeout:                hystrixTimeout,
		MaxConcurrentRequests:  hystrixMaxConcurrentRequests,
		RequestVolumeThreshold: hystrixRequestVolumeThreshold,
		SleepWindow:            hystrixSleepWindow,
		ErrorPercentThreshold:  hystrixErrorPercentThreshold,
	})

}
