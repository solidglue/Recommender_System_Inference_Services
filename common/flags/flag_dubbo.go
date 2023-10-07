package flags

import "flag"

var flagDubboInstance *flagDubbo

type flagDubbo struct {
	dubboServiceFile *string
}

//INFO: singleton instance
func init() {
	flagDubboInstance = new(flagDubbo)
}

func getFlagDubboInstance() *flagDubbo {
	return flagDubboInstance
}

// dubbo_serverconf
func (s *flagDubbo) setDubboServiceFile() {
	conf := flag.String("dubbo_serverconf", "conf/dubbogo_server.yml", "")
	s.dubboServiceFile = conf
}

func (s *flagDubbo) GetDubboServiceFile() *string {
	return s.dubboServiceFile
}
