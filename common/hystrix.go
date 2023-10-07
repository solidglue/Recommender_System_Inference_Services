package common

import (
	"infer-microservices/common/flags"

	"github.com/afex/hystrix-go/hystrix"
)

//https://zhuanlan.zhihu.com/p/540217601
//https://www.51cto.com/article/663588.html
//https://blog.csdn.net/skh2015java/article/details/110670390

// hystrix是加在客户端上的？

var hystrixTimeout int
var hystrixMaxConcurrentRequests int
var hystrixRequestVolumeThreshold int
var hystrixSleepWindow int
var hystrixErrorPercentThreshold int

func init() {

	//
	// timeout := *flags.Hystrix_timeoutMS
	// maxConcurrentRequests := *flags.Hystrix_MaxConcurrentRequests
	// requestVolumeThreshold := *flags.Hystrix_ErrorPercentThreshold
	// sleepWindow := *flags.Hystrix_SleepWindow
	// errorPercentThreshold := *flags.Hystrix_ErrorPercentThreshold

	flagFactory := flags.FlagFactory{}
	flagHystrix := flagFactory.FlagHystrixFactory()

	timeout := *flagHystrix.GetHystrixTimeoutMs()
	maxConcurrentRequests := *flagHystrix.GetHystrixMaxConcurrentRequests()
	requestVolumeThreshold := *flagHystrix.GetHystrixErrorPercentThreshold()
	sleepWindow := *flagHystrix.GetHystrixSleepWindow()
	errorPercentThreshold := *flagHystrix.GetHystrixErrorPercentThreshold()

	//all ConfigureCommands use the same conf, we can use different conf if we need.
	hystrix.ConfigureCommand("dubboServer", hystrix.CommandConfig{
		Timeout:                timeout,                // 执行 command 的超时时间
		MaxConcurrentRequests:  maxConcurrentRequests,  // 最大并发量,如果请求超过这个最大值将拒绝后续的请求，默认值为10
		RequestVolumeThreshold: requestVolumeThreshold, //  // 一个统计窗口 10 秒内请求数量 , 达到这个请求数量后才去判断是否要开启熔断
		SleepWindow:            sleepWindow,            // 	 // SleepWindow 的时间就是控制过多久后去尝试服务是否可用了  单位为毫秒
		ErrorPercentThreshold:  errorPercentThreshold,  // 错误百分比 请求数量大于等于 RequestVolumeThreshold 并且错误率到达这个百分比后就会启动熔断
	})

	hystrix.ConfigureCommand("restServer", hystrix.CommandConfig{
		Timeout:                timeout,                // 执行 command 的超时时间
		MaxConcurrentRequests:  maxConcurrentRequests,  // 最大并发量,如果请求超过这个最大值将拒绝后续的请求，默认值为10
		RequestVolumeThreshold: requestVolumeThreshold, //  // 一个统计窗口 10 秒内请求数量 , 达到这个请求数量后才去判断是否要开启熔断
		SleepWindow:            sleepWindow,            // 	 // SleepWindow 的时间就是控制过多久后去尝试服务是否可用了  单位为毫秒
		ErrorPercentThreshold:  errorPercentThreshold,  // 错误百分比 请求数量大于等于 RequestVolumeThreshold 并且错误率到达这个百分比后就会启动熔断
	})

	hystrix.ConfigureCommand("grpcServer", hystrix.CommandConfig{
		Timeout:                timeout,                // 执行 command 的超时时间
		MaxConcurrentRequests:  maxConcurrentRequests,  // 最大并发量,如果请求超过这个最大值将拒绝后续的请求，默认值为10
		RequestVolumeThreshold: requestVolumeThreshold, //  // 一个统计窗口 10 秒内请求数量 , 达到这个请求数量后才去判断是否要开启熔断
		SleepWindow:            sleepWindow,            // 	 // SleepWindow 的时间就是控制过多久后去尝试服务是否可用了  单位为毫秒
		ErrorPercentThreshold:  errorPercentThreshold,  // 错误百分比 请求数量大于等于 RequestVolumeThreshold 并且错误率到达这个百分比后就会启动熔断
	})

}
