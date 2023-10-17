package redis_config_loader

type RedisFactory struct {
}

func (r *RedisFactory) RedisClientFactory(domain string, dataId string, redisConfStr string) *RedisClient {
	rf := new(RedisClient)
	rf.ConfigLoad(domain, dataId, redisConfStr)

	return rf
}
