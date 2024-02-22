package redis

import (
	"context"
	"infer-microservices/internal/flags"
	"time"

	redis "github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var password string
var flagRedis flags.FlagRedis

type InferRedisClient struct {
	cli *redis.ClusterClient
}

func init() {
	flagFactory := flags.FlagFactory{}
	flagRedis := flagFactory.CreateFlagRedis()
	password = *flagRedis.GetRedisPassword()
}

func NewRedisClusterClient(data map[string]interface{}) *InferRedisClient {
	addrs_raw := data["addrs"].([]interface{})
	readTimeout := time.Duration(int64(data["readTimeoutMs"].(float64))) * time.Millisecond
	writeTimeout := time.Duration(int64(data["writeTimeoutMs"].(float64))) * time.Millisecond
	dialTimeout := time.Duration(int64(data["dialTimeoutMs"].(float64))) * time.Millisecond
	idleTimeout := time.Duration(int64(data["idleTimeoutS"].(float64))) * time.Second
	maxRetries := int(data["maxRetries"].(float64))
	minIdleConns := int(data["minIdleConns"].(float64))

	addrs := make([]string, 0)
	for _, addr := range addrs_raw {
		addrs = append(addrs, addr.(string))
	}
	cli := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:         addrs,
		Password:      password,
		ReadTimeout:   readTimeout,
		WriteTimeout:  writeTimeout,
		DialTimeout:   dialTimeout,
		IdleTimeout:   idleTimeout,
		PoolTimeout:   4 * time.Second,
		MaxRetries:    maxRetries,
		MinIdleConns:  minIdleConns,
		ReadOnly:      true,
		RouteRandomly: true,
	})

	return &InferRedisClient{cli: cli}
}

func (m *InferRedisClient) Get(key string) (string, error) {
	cmd := m.cli.Get(ctx, key)
	value, err := cmd.Result()
	if err != nil {
		return "", err
	}

	return value, nil
}

func (m *InferRedisClient) Set(key string, value string, expire time.Duration) error {
	return m.cli.Set(ctx, key, value, expire).Err()
}

func (m *InferRedisClient) HGet(key string, field string) (string, error) {
	value, err := m.cli.HGet(ctx, key, field).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

func (m *InferRedisClient) HSet(key string, field string, value string) error {
	return m.cli.HSet(ctx, key, field, value).Err()
}

func (m *InferRedisClient) HGetAll(key string) (map[string]string, error) {
	value, err := m.cli.HGetAll(ctx, key).Result()
	if err != nil {
		return map[string]string{}, err
	}
	return value, nil
}

// list
func (m *InferRedisClient) Lrange(key string) ([]string, error) {
	cmd := m.cli.LRange(ctx, key, 0, -1)
	value, err := cmd.Result()
	if err != nil {
		return make([]string, 0), err
	}

	return value, nil
}

func (m *InferRedisClient) Lrem(key string, value string) error {
	return m.cli.LRem(ctx, key, 0, value).Err()
}

func (m *InferRedisClient) Exists(key string) (bool, error) {
	ok, err := m.cli.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	if ok == 1 {
		return true, nil
	} else {
		return false, nil
	}
}

// INFO: Expire shoud after Push
func (m *InferRedisClient) LPush(key string, value string, expire time.Duration) error {
	m.cli.LPush(ctx, key, value)
	return m.cli.Expire(ctx, key, expire).Err()
}

func (m *InferRedisClient) RPush(key string, value string, expire time.Duration) error {
	m.cli.RPush(ctx, key, value)
	return m.cli.Expire(ctx, key, expire).Err()
}

// del
func (m *InferRedisClient) Delete(key string) error {
	return m.cli.Del(ctx, key).Err()
}
