package redis_config

import (
	"testing"
)

func TestRedisConfigLoader(t *testing.T) {

	testKey := "test_userid:1001"
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
	redisConf := RedisConfig{}
	redisConf.ConfigLoad("testId", redisTestStr)
	t.Log("redisConf:", redisConf)

	rst, err := redisConf.redisPool.HGetAll(testKey)
	if err != nil {
		t.Errorf("redis init failed")
	}

	t.Log("redis result:", rst)
}
