package redis_config_loader

import (
	"infer-microservices/utils"
	redis_v8 "infer-microservices/utils/redis"
)

type RedisConfig struct {
	domain    string //share redis by domain.
	redisPool *redis_v8.InferRedisClient
}

// // INFO: singleton instance
// func init() {
// 	RedisClientInstance = new(RedisClient)
// }

// func getRedisClientInstance() *RedisClient {
// 	return RedisClientInstance
// }

// domain
func (r *RedisConfig) setRedisDomain(domain string) {
	r.domain = domain
}

func (r *RedisConfig) GetRedisDomain() string {
	return r.domain
}

// redis pool
func (r *RedisConfig) setRedisPool(redisPool *redis_v8.InferRedisClient) {
	r.redisPool = redisPool
}

func (r *RedisConfig) GetRedisPool() *redis_v8.InferRedisClient {
	return r.redisPool
}

// @implement ConfigLoadInterface
func (r *RedisConfig) ConfigLoad(dataId string, redisConfStr string) error {
	confMap := utils.ConvertJsonToStruct(redisConfStr)
	redisClusterInfo := confMap["redisClusterName"].(map[string]interface{})
	redisConnPool := redis_v8.NewRedisClusterClient(redisClusterInfo)

	r.setRedisDomain(dataId)
	r.setRedisPool(redisConnPool)

	return nil
}
