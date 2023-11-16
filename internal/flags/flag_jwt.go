package flags

type FlagJwt struct {
	jwtKey *string
}

var flagJwtInstance *FlagJwt

// singleton instance
func init() {
	flagJwtInstance = new(FlagJwt)
}

func getFlagJwtInstance() *FlagJwt {
	return flagJwtInstance
}

// jwtKey
func (s *FlagJwt) setJwtKey(jwtKey *string) {
	s.jwtKey = jwtKey
}

func (s *FlagJwt) GetJwtKey() *string {
	return s.jwtKey
}
