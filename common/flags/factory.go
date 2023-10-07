package flags

import "flag"

type FlagFactory struct {
}

func init() {
	flag.Parse()
}

//start file factory
func (f *FlagFactory) FlagServiceConfigFactory() *FlagServiceStartInfo {
	fs := getFlagServiceStartInfoInstance()
	fs.setServiceConfigFile()
	fs.setServiceRestPort()
	fs.GetServiceGrpcPort()
	fs.setServiceMaxCpuNum()

	return fs
}

//cache factory
func (f *FlagFactory) FlagCacheFactory() *flagCache {
	fc := getFlagCacheInstance()
	fc.setBigcacheShards()
	fc.setBigcacheLifeWindowS()
	fc.setBigcacheCleanWindowS()
	fc.setBigcacheHardMaxCacheSize()
	fc.setBigcacheMaxEntrySize()
	fc.setBigcacheMaxEntriesInWindow()
	fc.setBigcacheVerbose()

	return fc
}

//dubbo factory
func (f *FlagFactory) FlagDubboFactory() *flagDubbo {
	fd := getFlagDubboInstance()
	fd.setDubboServiceFile()

	return fd
}

//dystrix factory
func (f *FlagFactory) FlagHystrixFactory() *flagsHystrix {
	fh := getFlagsHystrixInstance()
	fh.setHystrixErrorPercentThreshold()
	fh.setHystrixLowerRankNum()
	fh.setHystrixLowerRecallNum()
	fh.setHystrixMaxConcurrentRequests()
	fh.setHystrixRequestVolumeThreshold()
	fh.setHystrixSleepWindow()
	fh.setHystrixTimeoutMs()

	return fh
}

//logs factory
func (f *FlagFactory) FlagLogFactory() *flagsLog {
	fl := getFlagLogInstance()
	fl.setLogFileName()
	fl.setLogLevel()
	fl.setLogMaxSize()
	fl.setLogSaveDays()

	return fl
}

//nacos factory
func (f *FlagFactory) FlagNacosFactory() *flagsNacos {
	fn := getFlagsNacosInstance()
	fn.setNacosIp()
	fn.setNacosPort()
	fn.setNacosUsername()
	fn.setNacosPassword()
	fn.setNacosLogdir()
	fn.setNacosLoglevel()
	fn.setNacosCachedir()
	fn.setNacosTimeoutMs()

	return fn
}

//redis factory
func (f *FlagFactory) FlagRedisFactory() *FlagRedis {
	fd := getFlagRedisInstance()
	fd.setRedisPassword()

	return fd
}

//skywalking factory
func (f *FlagFactory) FlagSkywalkingFactory() *flagsSkywalking {
	fs := getFlagsSkywalkingInstance()
	fs.setSkywalkingWhetheropen()
	fs.setSkywalkingIp()
	fs.setSkywalkingPort()
	fs.setSkywalkingServername()

	return fs
}

//tensorflow factory
func (f *FlagFactory) FlagTensorflowFactory() *flagTensorflow {
	ft := getFlagTensorflowInstance()
	ft.setTfservingModelVersion()
	ft.setTfservingTimeoutMs()

	return ft
}
