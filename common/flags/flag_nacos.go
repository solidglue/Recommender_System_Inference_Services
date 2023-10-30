package flags

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
func (s *flagsNacos) setNacosIp(nacosIp *string) {
	s.nacosIp = nacosIp
}

func (s *flagsNacos) GetNacosIp() *string {
	return s.nacosIp
}

// nacos_port
func (s *flagsNacos) setNacosPort(nacosPort *int) {
	s.nacosPort = nacosPort
}

func (s *flagsNacos) GetNacosPort() *int {
	return s.nacosPort
}

// nacos_username
func (s *flagsNacos) setNacosUsername(nacosUsername *string) {
	s.nacosUsername = nacosUsername
}

func (s *flagsNacos) GetNacosUsername() *string {
	return s.nacosUsername
}

// nacos_password
func (s *flagsNacos) setNacosPassword(nacosPassword *string) {
	s.nacosPassword = nacosPassword
}

func (s *flagsNacos) GetNacosPassword() *string {
	return s.nacosPassword
}

// nacos_logdir
func (s *flagsNacos) setNacosLogdir(nacosLogdir *string) {
	s.nacosLogdir = nacosLogdir
}

func (s *flagsNacos) GetNacosLogdir() *string {
	return s.nacosLogdir
}

// nacos_cachedir
func (s *flagsNacos) setNacosCachedir(nacosCachedir *string) {
	s.nacosCachedir = nacosCachedir
}

func (s *flagsNacos) GetacosCachedir() *string {
	return s.nacosCachedir
}

// nacos_loglevel
func (s *flagsNacos) setNacosLoglevel(nacosLoglevel *string) {
	s.nacosLoglevel = nacosLoglevel
}

func (s *flagsNacos) GetNacosLoglevel() *string {
	return s.nacosLoglevel
}

// nacos_timeoutMS
func (s *flagsNacos) setNacosTimeoutMs(nacosTimeoutMs *int) {
	s.nacosTimeoutMs = nacosTimeoutMs
}

func (s *flagsNacos) GetNacosTimeoutMs() *int {
	return s.nacosTimeoutMs
}
