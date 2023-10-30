package flags

var flagsSkywalkingInstance *flagsSkywalking

type flagsSkywalking struct {
	skywalkingWhetheropen *bool
	skywalkingServername  *string
	skywalkingIp          *string
	skywalkingPort        *int
}

// singleton instance
func init() {
	flagsSkywalkingInstance = new(flagsSkywalking)
}

func getFlagsSkywalkingInstance() *flagsSkywalking {
	return flagsSkywalkingInstance
}

// skywalking_whetheropen
func (s *flagsSkywalking) setSkywalkingWhetheropen(skywalkingWhetheropen *bool) {
	s.skywalkingWhetheropen = skywalkingWhetheropen
}

func (s *flagsSkywalking) GetSkywalkingWhetheropen() *bool {
	return s.skywalkingWhetheropen
}

// skywalking_servername
func (s *flagsSkywalking) setSkywalkingServername(skywalkingServername *string) {
	s.skywalkingServername = skywalkingServername
}

func (s *flagsSkywalking) GetSkywalkingServername() *string {
	return s.skywalkingServername
}

// skywalking_ip
func (s *flagsSkywalking) setSkywalkingIp(skywalkingIp *string) {
	s.skywalkingIp = skywalkingIp
}

func (s *flagsSkywalking) GetSkywalkingIp() *string {
	return s.skywalkingIp
}

// skywalking_port
func (s *flagsSkywalking) setSkywalkingPort(skywalkingPort *int) {
	s.skywalkingPort = skywalkingPort
}

func (s *flagsSkywalking) GetSkywalkingPort() *int {
	return s.skywalkingPort
}
