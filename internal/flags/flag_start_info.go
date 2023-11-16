package flags

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
func (s *FlagServiceStartInfo) setServiceConfigFile(serviceConfig *string) {
	s.serviceConfig = serviceConfig
}

func (s *FlagServiceStartInfo) GetServiceConfigFile() *string {
	return s.serviceConfig
}

// rest_server_port
func (s *FlagServiceStartInfo) setServiceRestPort(restPort *uint) {
	s.restPort = restPort
}

func (s *FlagServiceStartInfo) GetServiceRestPort() *uint {
	return s.restPort
}

// grpc_server_port
func (s *FlagServiceStartInfo) setServiceGrpcPort(grpcPort *uint) {
	s.grpcPort = grpcPort
}

func (s *FlagServiceStartInfo) GetServiceGrpcPort() *uint {
	return s.grpcPort
}

// max_cpu_num
func (s *FlagServiceStartInfo) setServiceMaxCpuNum(maxCpuNum *int) {
	s.maxCpuNum = maxCpuNum
}

func (s *FlagServiceStartInfo) GetServiceMaxCpuNum() *int {
	return s.maxCpuNum
}
