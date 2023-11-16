package redis_config

import (
	redis_v8 "infer-microservices/pkg/db/redis"
	"infer-microservices/pkg/utils"
)

type RedisConfig struct {
	redisPool *redis_v8.InferRedisClient `validate:"required"`
}

// // INFO: singleton instance
// func init() {
// 	RedisClientInstance = new(RedisClient)
// }

// func getRedisClientInstance() *RedisClient {
// 	return RedisClientInstance
// }

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

	r.setRedisPool(redisConnPool)

	return nil
}
