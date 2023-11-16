package redis

import (
	"fmt"
	"strings"
	"time"

	redis "github.com/go-redis/redis/v8"
)

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
