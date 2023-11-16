package flags

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
func (s *FlagRedis) setRedisPassword(redisPassword *string) {
	s.redisPassword = redisPassword
}

func (s *FlagRedis) GetRedisPassword() *string {
	return s.redisPassword
}
