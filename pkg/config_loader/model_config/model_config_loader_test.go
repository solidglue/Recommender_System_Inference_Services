package model_config

import (
	"infer-microservices/pkg/config_loader/redis_config"
	"testing"
)

func TestModelConfigLoader(t *testing.T) {

	testUserId := "u1111"
	testItemId := "i7001"

	redisTestStr := `
	{
		"redis_conf": {		
			"redisCluster": {
				"addrs": [],
				"password": "",
				"idleTimeoutMs": 100,
				"readTimeoutMs": 100,
				"writeTimeoutMs":100,
				"dialTimeoutS":600,
				"maxRetries":2,
				"minIdleConns":50
			}
		}
	}
	`

	modelTestStr := `
	{
		"model_conf": {
			"model-001": {
				"fieldsSpec": [{},{},{}],
				"tfservingGrpcAddr": {
					"tfservingModelName": "models",
					"addrs": [],
					"pool_size": 50,
					"initCap":10,
					"idleTimeoutMs": 100,
					"readTimeoutMs": 100,
					"writeTimeoutMs":100,
					"dialTimeoutS":600
				},
				"user_feature_rediskey_pre": "",
				"item_feature_rediskey_pre": ""
			}
		},
	}
	`
	redisConf := redis_config.RedisConfig{}
	redisConf.ConfigLoad("testId", redisTestStr)

	modelConf := ModelConfig{}
	modelConf.ConfigLoad("testId", modelTestStr)
	t.Log("modelConf:", modelConf)

	//test user features
	rst, err := redisConf.GetRedisPool().HGetAll(modelConf.userRedisKeyPre + testUserId)
	if err != nil {
		t.Errorf("redis init failed")
	}
	t.Log("user features:", rst)

	//test item features
	rst, err = redisConf.GetRedisPool().HGetAll(modelConf.userRedisKeyPre + testItemId)
	if err != nil {
		t.Errorf("redis init failed")
	}
	t.Log("item features:", rst)
}
