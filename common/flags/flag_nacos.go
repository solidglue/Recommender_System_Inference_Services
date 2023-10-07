package flags

import "flag"

var flagsNacosInstance *flagsNacos

type flagsNacos struct {

	//nacos
	nacosIp        *string
	nacosPort      *int
	nacosUsername  *string
	nacosPassword  *string
	nacosLogdir    *string
	nacosCachedir  *string
	nacosLoglevel  *string
	nacosTimeoutMs *int
}

// singleton instance
func init() {
	flagsNacosInstance = new(flagsNacos)
}

func getFlagsNacosInstance() *flagsNacos {
	return flagsNacosInstance
}

// nacos_ip
func (s *flagsNacos) setNacosIp() {
	conf := flag.String("nacos_ip", "10.10.10.10", "")
	s.nacosIp = conf
}

func (s *flagsNacos) GetNacosIp() *string {
	return s.nacosIp
}

// nacos_port
func (s *flagsNacos) setNacosPort() {
	conf := flag.Int("nacos_port", 8081, "")
	s.nacosPort = conf
}

func (s *flagsNacos) GetNacosPort() *int {
	return s.nacosPort
}

// nacos_username
func (s *flagsNacos) setNacosUsername() {
	conf := flag.String("nacos_username", "nacos", "")
	s.nacosUsername = conf
}

func (s *flagsNacos) GetNacosUsername() *string {
	return s.nacosUsername
}

// nacos_password
func (s *flagsNacos) setNacosPassword() {
	conf := flag.String("nacos_password", "nacos", "")
	s.nacosPassword = conf
}

func (s *flagsNacos) GetNacosPassword() *string {
	return s.nacosPassword
}

// nacos_logdir
func (s *flagsNacos) setNacosLogdir() {
	conf := flag.String("nacos_logdir", "nacos-logs", "")
	s.nacosLogdir = conf
}

func (s *flagsNacos) GetNacosLogdir() *string {
	return s.nacosLogdir
}

// nacos_cachedir
func (s *flagsNacos) setNacosCachedir() {
	conf := flag.String("nacos_cachedir", "nacos-cache", "")
	s.nacosCachedir = conf
}

func (s *flagsNacos) GetacosCachedir() *string {
	return s.nacosCachedir
}

// nacos_loglevel
func (s *flagsNacos) setNacosLoglevel() {
	conf := flag.String("nacos_loglevel", "error", "")
	s.nacosLoglevel = conf
}

func (s *flagsNacos) GetNacosLoglevel() *string {
	return s.nacosLoglevel
}

// nacos_timeoutMS
func (s *flagsNacos) setNacosTimeoutMs() {
	conf := flag.Int("nacos_timeoutMS", 5000, "")
	s.nacosTimeoutMs = conf
}

func (s *flagsNacos) GetNacosTimeoutMs() *int {
	return s.nacosTimeoutMs
}
