package flags

import (
	"flag"
)

var flagServiceStartInfoInstance *FlagServiceStartInfo

type FlagServiceStartInfo struct {
	serviceConfig *string
	restPort      *uint
	grpcPort      *uint
	maxCpuNum     *int
}

// singleton instance
func init() {
	flagServiceStartInfoInstance = new(FlagServiceStartInfo)
}

func getFlagServiceStartInfoInstance() *FlagServiceStartInfo {
	return flagServiceStartInfoInstance
}

// service_start_file
func (s *FlagServiceStartInfo) setServiceConfigFile() {
	conf := flag.String("service_start_file", "./conf/service_conf_file.json", "")
	s.serviceConfig = conf
}

func (s *FlagServiceStartInfo) GetServiceConfigFile() *string {
	return s.serviceConfig
}

// rest_server_port
func (s *FlagServiceStartInfo) setServiceRestPort() {
	conf := flag.Uint("rest_server_port", 8888, "")
	s.restPort = conf
}

func (s *FlagServiceStartInfo) GetServiceRestPort() *uint {
	return s.restPort
}

// grpc_server_port
func (s *FlagServiceStartInfo) setServiceGrpcPort() {
	conf := flag.Uint("grpc_server_port", 8889, "")
	s.grpcPort = conf
}

func (s *FlagServiceStartInfo) GetServiceGrpcPort() *uint {
	return s.grpcPort
}

// max_cpu_num
func (s *FlagServiceStartInfo) setServiceMaxCpuNum() {
	conf := flag.Int("max_cpu_num", 16, "")
	s.maxCpuNum = conf
}

func (s *FlagServiceStartInfo) GetServiceMaxCpuNum() *int {
	return s.maxCpuNum
}
