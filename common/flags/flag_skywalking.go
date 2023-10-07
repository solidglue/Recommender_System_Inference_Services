package flags

import "flag"

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
func (s *flagsSkywalking) setSkywalkingWhetheropen() {
	conf := flag.Bool("skywalking_whetheropen", false, "")
	s.skywalkingWhetheropen = conf
}

func (s *flagsSkywalking) GetSkywalkingWhetheropen() *bool {
	return s.skywalkingWhetheropen
}

// skywalking_servername
func (s *flagsSkywalking) setSkywalkingServername() {
	conf := flag.String("skywalking_servername", "infer", "")
	s.skywalkingServername = conf
}

func (s *flagsSkywalking) GetSkywalkingServername() *string {
	return s.skywalkingServername
}

// skywalking_ip
func (s *flagsSkywalking) setSkywalkingIp() {
	conf := flag.String("skywalking_ip", "10.10.10.10", "")
	s.skywalkingIp = conf
}

func (s *flagsSkywalking) GetSkywalkingIp() *string {
	return s.skywalkingIp
}

// skywalking_port
func (s *flagsSkywalking) setSkywalkingPort() {
	conf := flag.Int("skywalking_port", 8080, "")
	s.skywalkingPort = conf
}

func (s *flagsSkywalking) GetSkywalkingPort() *int {
	return s.skywalkingPort
}
