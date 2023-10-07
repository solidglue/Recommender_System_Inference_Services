package common

import (
	"fmt"
	"infer-microservices/common/flags"
	"strings"
	"time"

	//"github.com/go-redis/redis"
	"context"

	redis "github.com/go-redis/redis/v8"
)

//TODO:加密访问redis，已删除
//TODO:redis pipeline  https://zhuanlan.zhihu.com/p/613192456

//Pipeline 主要是一种网络优化,它本质上意味着客户端缓冲一堆命令并一次性将它们发送到服务器，减少了每条命令分别传输的IO开销, 同时减少了系统调用的次数，因此提升了整体的吞吐能力,节省了每个命令的网络往返时间（RTT）。

var ctx = context.Background()
var password string
var flagRedis flags.FlagRedis

type InferRedisClient struct {
	cli *redis.ClusterClient
}

func init() {

	flagFactory := flags.FlagFactory{}
	flagRedis := flagFactory.FlagRedisFactory()
	password = *flagRedis.GetRedisPassword() //密码秘钥分开存储

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

	//https://www.cnblogs.com/CJ-cooper/p/15149273.html    //https://www.lixueduan.com/posts/redis/db-connection-pool-settings/
	cli := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:         addrs,
		Password:      password,
		ReadTimeout:   readTimeout,
		WriteTimeout:  writeTimeout,
		DialTimeout:   dialTimeout,
		IdleTimeout:   idleTimeout,
		PoolTimeout:   4 * time.Second, //add
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

func (m *InferRedisClient) HGetAll(key string) (map[string]string, error) {

	value, err := m.cli.HGetAll(ctx, key).Result()
	if err != nil {
		return map[string]string{}, err
	}
	return value, nil
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

func (m *InferRedisClient) GetPipe(keys *[]string, rate int) (map[string]string, error) {

	res := make(map[string]string)
	cmders := []redis.Cmder{}
	pipe := m.cli.Pipeline()

	for i, key := range *keys {

		pipe.Get(ctx, key)

		if i > 0 && i%rate == 0 {
			perResult, err := pipe.Exec(ctx)
			if err != nil {
				return res, err
			}
			cmders = append(cmders, perResult...)
		}
	}
	perResult, err := pipe.Exec(ctx)
	if err != nil {

		return res, err
	}

	cmders = append(cmders, perResult...)

	valueArr := []string{}
	for _, cmder := range cmders {

		cmd := cmder.(*redis.StringCmd)
		str, err := cmd.Result()

		if err != nil {
			return res, err
		}

		valueArr = append(valueArr, str)
	}

	for i, key := range *keys {
		res[key] = valueArr[i]
	}

	return res, nil
}

func (m *InferRedisClient) SetPipe(kv *map[string]string, rate int, expire time.Duration) error {

	pipe := m.cli.Pipeline()
	i := 0

	for k, v := range *kv {
		pipe.Set(ctx, k, v, expire)
		i++
		if i%rate == 0 {
			_, err := pipe.Exec(ctx)
			if err != nil {
				return err
			}
		}
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (m *InferRedisClient) HGetPipe(kf *map[string]interface{}, rate int) (map[string]interface{}, error) {

	res := make(map[string]interface{})
	cmders := []redis.Cmder{}
	pipe := m.cli.Pipeline()
	i := 0
	kfArr := []string{}

	for key, fields := range *kf {

		fieldsArr := fields.([]string)

		for _, field := range fieldsArr {

			i++
			kfArr = append(kfArr, fmt.Sprintf("%s###%s", key, field))
			pipe.HGet(ctx, key, field)

			if i%rate == 0 {
				perResult, err := pipe.Exec(ctx)
				if err != nil {
					return res, err
				}
				cmders = append(cmders, perResult...)
			}
		}
	}
	perResult, err := pipe.Exec(ctx)
	if err != nil {
		return res, err
	}
	cmders = append(cmders, perResult...)

	valueArr := []string{}
	for _, cmder := range cmders {

		cmd := cmder.(*redis.StringCmd)
		str, err := cmd.Result()
		if err != nil {
			return res, err
		}

		valueArr = append(valueArr, str)
	}

	for _, kf := range kfArr {
		println(kf)
	}

	for _, value := range valueArr {
		println(value)
	}

	for i, kf := range kfArr {

		key := strings.Split(kf, "###")[0]
		field := strings.Split(kf, "###")[1]

		if _, ok := res[key]; ok {

			fv := res[key].(map[string]string)
			fv[field] = valueArr[i]
			res[key] = fv
		} else {

			fv := make(map[string]string)
			fv[field] = valueArr[i]
			res[key] = fv

		}

	}

	return res, nil

}

func (m *InferRedisClient) HSetPipe(kfv *map[string]interface{}, rate int) error {

	pipe := m.cli.Pipeline()
	i := 0
	for key, fv := range *kfv {

		for field, value := range fv.(map[string]string) {

			i++
			pipe.HSet(ctx, key, field, value)

			if i%rate == 0 {

				_, err := pipe.Exec(ctx)
				if err != nil {
					return err
				}
			}

		}
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (m *InferRedisClient) PipeExist(keys *[]string, rate int) (map[string]int64, error) {

	startTime := time.Now()

	cmders := []redis.Cmder{}
	map_result := make(map[string]int64)
	pipe := m.cli.Pipeline()

	defer pipe.Close()

	for i, key := range *keys {

		pipe.Exists(ctx, key)
		if (i+1)%rate == 0 {
			perResult, err := pipe.Exec(ctx)
			if err != nil {
				return map_result, err
			}
			cmders = append(cmders, perResult...)
		}
	}

	elapsed := time.Since(startTime)
	fmt.Println("elapsed=", elapsed)

	startTime = time.Now()

	perResult, err := pipe.Exec(ctx)
	if err != nil {
		return map_result, err
	}
	cmders = append(cmders, perResult...)
	elapsed = time.Since(startTime)
	fmt.Println("elapsed=", elapsed)

	startTime = time.Now()

	result := make([]int64, 0)
	for _, cmder := range cmders {

		cmd := cmder.(*redis.IntCmd)
		f, err := cmd.Result()
		if err != nil {
			return map_result, err
		}

		result = append(result, f)
	}

	for i, key := range *keys {
		map_result[key] = result[i]
	}

	elapsed = time.Since(startTime)
	fmt.Println("elapsed=", elapsed)

	return map_result, nil
}

func (m *InferRedisClient) HGetPipeAll(keys *[]string, rate int) (map[string]interface{}, error) {

	res := make(map[string]interface{})
	cmders := []redis.Cmder{}

	pipe := m.cli.Pipeline()

	defer pipe.Close()

	for i, key := range *keys {

		pipe.HGetAll(ctx, key)
		if (i+1)%rate == 0 {

			perResult, err := pipe.Exec(ctx)

			if err != nil {
				return res, nil
			}
			cmders = append(cmders, perResult...)
		}
	}

	perResult, err := pipe.Exec(ctx)

	if err != nil {
		return res, err
	}
	cmders = append(cmders, perResult...)

	valueArr := make([]map[string]string, 0)

	for _, cmder := range cmders {

		cmd := cmder.(*redis.StringStringMapCmd)
		mss, err := cmd.Result()
		if err != nil {
			return res, err
		}
		valueArr = append(valueArr, mss)
	}

	for i, key := range *keys {

		res[key] = valueArr[i]

	}

	return res, nil
}
