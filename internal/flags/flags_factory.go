package flags

import "flag"

type FlagFactory struct {
}

func init() {
	flag.Parse()
}

//start file factory
func (f *FlagFactory) CreateFlagServiceConfig() *FlagServiceStartInfo {
	serviceStartFile := flag.String("service_start_file", "./conf/service_conf_file.json", "")
	restServerPort := flag.Uint("rest_server_port", 8888, "")
	grpcServerPort := flag.Uint("grpc_server_port", 8889, "")
	maxCpuNum := flag.Int("max_cpu_num", 16, "")

	fs := getFlagServiceStartInfoInstance()
	fs.setServiceConfigFile(serviceStartFile)
	fs.setServiceRestPort(restServerPort)
	fs.setServiceGrpcPort(grpcServerPort)
	fs.setServiceMaxCpuNum(maxCpuNum)

	return fs
}

//cache factory
func (f *FlagFactory) CreateFlagCache() *flagCache {
	bigcaheShards := flag.Int("bigcahe_shards", 1024, "")
	bigcaheLifeWindowS := flag.Int("bigcahe_lifeWindowS", 300, "")
	bigcacheCleanWindowS := flag.Int("bigcache_cleanWindowS", 120, "")
	bigcacheHardMaxCacheSize := flag.Int("bigcache_hardMaxCacheSize", 409600, "MB")
	bigcacheMaxEntrySize := flag.Int("bigcache_maxEntrySize", 1024, "byte")
	bigcacheMaxEntriesInWindow := flag.Int("bigcache_maxEntriesInWindow", 2000000, "depends on tps")
	bigcacheVerbose := flag.Bool("bigcache_verbose", false, "")

	fc := getFlagCacheInstance()
	fc.setBigcacheShards(bigcaheShards)
	fc.setBigcacheLifeWindowS(bigcaheLifeWindowS)
	fc.setBigcacheCleanWindowS(bigcacheCleanWindowS)
	fc.setBigcacheHardMaxCacheSize(bigcacheHardMaxCacheSize)
	fc.setBigcacheMaxEntrySize(bigcacheMaxEntrySize)
	fc.setBigcacheMaxEntriesInWindow(bigcacheMaxEntriesInWindow)
	fc.setBigcacheVerbose(bigcacheVerbose)

	return fc
}

//dubbo factory
func (f *FlagFactory) CreateFlagDubbo() *flagDubbo {
	fd := getFlagDubboInstance()
	dubboServerconf := flag.String("dubbo_serverconf", "conf/dubbogo_server.yml", "")
	fd.setDubboServiceFile(dubboServerconf)

	return fd
}

//dystrix factory
func (f *FlagFactory) CreateFlagHystrix() *flagsHystrix {
	hystrixTimeoutMS := flag.Int("hystrix_timeoutMS", 100, "")
	hystrixRequestVolumeThreshold := flag.Int("hystrix_RequestVolumeThreshold", 50000, "")
	hystrixSleepWindow := flag.Int("hystrix_SleepWindow", 10000, "")
	hystrixErrorPercentThreshold := flag.Int("hystrix_ErrorPercentThreshold", 1, "")
	hystrixLowerRecallNum := flag.Int("hystrix_lowerRecallNum", 100, "")
	hystrixLowerRankNum := flag.Int("hystrix_lowerRankNum", 100, "")
	hystrixMaxConcurrentRequests := flag.Int("hystrix_MaxConcurrentRequests", 10000, "")

	fh := getFlagsHystrixInstance()
	fh.setHystrixErrorPercentThreshold(hystrixErrorPercentThreshold)
	fh.setHystrixLowerRankNum(hystrixLowerRankNum)
	fh.setHystrixLowerRecallNum(hystrixLowerRecallNum)
	fh.setHystrixMaxConcurrentRequests(hystrixMaxConcurrentRequests)
	fh.setHystrixRequestVolumeThreshold(hystrixRequestVolumeThreshold)
	fh.setHystrixSleepWindow(hystrixSleepWindow)
	fh.setHystrixTimeoutMs(hystrixTimeoutMS)

	return fh
}

//logs factory
func (f *FlagFactory) CreateFlagLog() *flagsLog {
	logMaxSize := flag.Int("log_max_size", 200000000, "the max size of the log file (in Byte)")
	logSaveDays := flag.Int("log_save_days", 7, "")
	logFileName := flag.String("log_file_name", "infer.log", "")
	logLevel := flag.String("log_level", "error", "the log level, (debug, info, error, fatal)")

	fl := getFlagLogInstance()
	fl.setLogFileName(logFileName)
	fl.setLogLevel(logLevel)
	fl.setLogMaxSize(logMaxSize)
	fl.setLogSaveDays(logSaveDays)

	return fl
}

//nacos factory
func (f *FlagFactory) CreateFlagNacos() *flagsNacos {
	nacosIp := flag.String("nacos_ip", "10.10.10.10", "")
	nacosPort := flag.Int("nacos_port", 8888, "")
	nacosUsername := flag.String("nacos_username", "nacos", "")
	nacosPassword := flag.String("nacos_password", "nacos", "")
	nacosLogdir := flag.String("nacos_logdir", "nacos-logs", "")
	nacosCachedir := flag.String("nacos_cachedir", "nacos-cache", "")
	nacosLoglevel := flag.String("nacos_loglevel", "error", "")
	nacosTimeoutMS := flag.Int("nacos_timeoutMS", 5000, "")

	fn := getFlagsNacosInstance()
	fn.setNacosIp(nacosIp)
	fn.setNacosPort(nacosPort)
	fn.setNacosUsername(nacosUsername)
	fn.setNacosPassword(nacosPassword)
	fn.setNacosLogdir(nacosLogdir)
	fn.setNacosLoglevel(nacosLoglevel)
	fn.setNacosCachedir(nacosCachedir)
	fn.setNacosTimeoutMs(nacosTimeoutMS)

	return fn
}

//redis factory
func (f *FlagFactory) CreateFlagRedis() *FlagRedis {
	redisPassword := flag.String("redis_password", "", "")

	fd := getFlagRedisInstance()
	fd.setRedisPassword(redisPassword)

	return fd
}

//skywalking factory
func (f *FlagFactory) CreateFlagSkywalking() *flagsSkywalking {
	skywalkingWhetheropen := flag.Bool("skywalking_whetheropen", false, "")
	skywalkingServername := flag.String("skywalking_servername", "infer", "")
	skywalkingIp := flag.String("skywalking_ip", "10.10.10.10", "")
	skywalkingPort := flag.Int("skywalking_port", 8080, "")

	fs := getFlagsSkywalkingInstance()
	fs.setSkywalkingWhetheropen(skywalkingWhetheropen)
	fs.setSkywalkingIp(skywalkingIp)
	fs.setSkywalkingPort(skywalkingPort)
	fs.setSkywalkingServername(skywalkingServername)

	return fs
}

//tensorflow factory
func (f *FlagFactory) CreateFlagTensorflow() *flagTensorflow {
	tfservingModeVersion := flag.Int64("tfserving_model_version", 0, "")
	tfservingTimeoutms := flag.Int64("tfserving_timeoutms", 100, "")

	ft := getFlagTensorflowInstance()
	ft.setTfservingModelVersion(tfservingModeVersion)
	ft.setTfservingTimeoutMs(tfservingTimeoutms)

	return ft
}

//bloom factory
func (f *FlagFactory) CreateFlagBloom() *FlagBloom {
	userCountLevel := flag.Uint("user_count_level", 100000000, "")
	itemCountLevel := flag.Uint("item_count_level", 10000000, "")

	ft := getFlagBloomInstance()
	ft.setUserCountLevel(userCountLevel)
	ft.setItemCountLevel(itemCountLevel)

	return ft
}

//viper factory
func (f *FlagFactory) CreateFlagViper() *FlagViper {
	configName := flag.String("config_name", "bloom_filter", "")
	configType := flag.String("config_type", "json", "")
	configPath := flag.String("configName", "./conf/", "")

	ft := getFlagViperInstance()
	ft.setConfigName(configName)
	ft.setConfigType(configType)
	ft.setConfigPath(configPath)

	return ft
}

//jwt factory
func (f *FlagFactory) CreateFlagJwt() *FlagJwt {
	jwtKey := flag.String("jwt_key", "im your dad", "")

	ft := getFlagJwtInstance()
	ft.setJwtKey(jwtKey)

	return ft
}

//kafka factory
func (f *FlagFactory) CreateFlagKafka() *flagKafka {
	kafkaUrl := flag.String("kafka_url", "l27.0.0.1:9092", "")
	kafkaTopic := flag.String("kafka_topic", "kafka_topic_001", "")
	kafkaGroup := flag.String("kafka_group", "kafka_group_001", "")

	ft := getFlagKafkaInstance()
	ft.setKafkaUrl(*kafkaUrl)
	ft.setKafkaTopic(*kafkaTopic)
	ft.setKafkaGroup(*kafkaGroup)

	return ft
}
