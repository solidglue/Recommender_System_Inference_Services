package flags

type FlagBloom struct {
	//redis
	userCountLevel *uint
	itemCountLevel *uint
}

var flagBloomInstance *FlagBloom

// singleton instance
func init() {
	flagBloomInstance = new(FlagBloom)
}

func getFlagBloomInstance() *FlagBloom {
	return flagBloomInstance
}

// userCountLevel
func (s *FlagBloom) setUserCountLevel(userCountLevel *uint) {
	s.userCountLevel = userCountLevel
}

func (s *FlagBloom) GetUserCountLevel() *uint {
	return s.userCountLevel
}

// userCountLevel
func (s *FlagBloom) setItemCountLevel(userCountLevel *uint) {
	s.userCountLevel = userCountLevel
}

func (s *FlagBloom) GetItemCountLevel() *uint {
	return s.userCountLevel
}
