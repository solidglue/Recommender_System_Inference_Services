package flags

import "flag"

var flagRedisInstance *FlagRedis

type FlagRedis struct {
	//redis
	redisPassword *string
}

// singleton instance
func init() {
	flagRedisInstance = new(FlagRedis)
}

func getFlagRedisInstance() *FlagRedis {
	return flagRedisInstance
}

// redis_password
func (s *FlagRedis) setRedisPassword() {
	conf := flag.String("redis_password", "", "")
	s.redisPassword = conf
}

func (s *FlagRedis) GetRedisPassword() *string {
	return s.redisPassword
}
