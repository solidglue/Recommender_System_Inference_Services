package redis_config

import (
	redis_helper "infer-microservices/common/redis"
	redis_v8 "infer-microservices/common/redis"
	"infer-microservices/utils"
)

var RedisClientInstance *RedisClient

type RedisClient struct {
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
func (r *RedisClient) setRedisDomain(domain string) {
	r.domain = domain
}

func (r *RedisClient) GetRedisDomain() string {
	return r.domain
}

// redis pool
func (r *RedisClient) setRedisPool(redisPool *redis_v8.InferRedisClient) {
	r.redisPool = redisPool
}

func (r *RedisClient) GetRedisPool() *redis_v8.InferRedisClient {
	return r.redisPool
}

// redis conf load
func (r *RedisClient) ConfigLoad(domain string, dataId string, redisConfStr string) error {

	confMap := utils.Json2Map(redisConfStr)
	redisClusterInfo := confMap["redisCluster"].(map[string]interface{})
	redisConnPool := redis_helper.NewRedisClusterClient(redisClusterInfo)

	r.setRedisDomain(domain)
	r.setRedisPool(redisConnPool)

	return nil

}
