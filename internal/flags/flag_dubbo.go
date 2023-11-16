package flags

var flagDubboInstance *flagDubbo

type flagDubbo struct {
	dubboServiceFile *string
}

// INFO: singleton instance
func init() {
	flagDubboInstance = new(flagDubbo)
}

func getFlagDubboInstance() *flagDubbo {
	return flagDubboInstance
}

// dubbo_serverconf
func (s *flagDubbo) setDubboServiceFile(dubboServiceFile *string) {
	s.dubboServiceFile = dubboServiceFile
}

func (s *flagDubbo) GetDubboServiceFile() *string {
	return s.dubboServiceFile
}
