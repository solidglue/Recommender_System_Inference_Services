package flags

var flagsHystrixInstance *flagsHystrix

type flagsHystrix struct {
	// hystrix
	hystrixTimeoutMs              *int
	hystrixMaxConcurrentRequests  *int
	hystrixRequestVolumeThreshold *int
	hystrixSleepWindow            *int
	hystrixErrorPercentThreshold  *int
	hystrixLowerRecallNum         *int
	hystrixLowerRankNum           *int
}

// singleton instance
func init() {
	flagsHystrixInstance = new(flagsHystrix)
}

func getFlagsHystrixInstance() *flagsHystrix {
	return flagsHystrixInstance
}

// hystrix_timeoutMS
func (s *flagsHystrix) setHystrixTimeoutMs(hystrixTimeoutMs *int) {
	s.hystrixTimeoutMs = hystrixTimeoutMs
}

func (s *flagsHystrix) GetHystrixTimeoutMs() *int {
	return s.hystrixTimeoutMs
}

// hystrix_MaxConcurrentRequests
func (s *flagsHystrix) setHystrixMaxConcurrentRequests(hystrixMaxConcurrentRequests *int) {
	s.hystrixMaxConcurrentRequests = hystrixMaxConcurrentRequests
}

func (s *flagsHystrix) GetHystrixMaxConcurrentRequests() *int {
	return s.hystrixMaxConcurrentRequests
}

// hystrix_RequestVolumeThreshold
func (s *flagsHystrix) setHystrixRequestVolumeThreshold(hystrixRequestVolumeThreshold *int) {
	s.hystrixRequestVolumeThreshold = hystrixRequestVolumeThreshold
}

func (s *flagsHystrix) GetHystrixRequestVolumeThreshold() *int {
	return s.hystrixRequestVolumeThreshold
}

// hystrix_timeoutMS
func (s *flagsHystrix) setHystrixSleepWindow(hystrixSleepWindow *int) {
	s.hystrixSleepWindow = hystrixSleepWindow
}

func (s *flagsHystrix) GetHystrixSleepWindow() *int {
	return s.hystrixSleepWindow
}

// hystrix_ErrorPercentThreshold
func (s *flagsHystrix) setHystrixErrorPercentThreshold(hystrixErrorPercentThreshold *int) {
	s.hystrixErrorPercentThreshold = hystrixErrorPercentThreshold
}

func (s *flagsHystrix) GetHystrixErrorPercentThreshold() *int {
	return s.hystrixErrorPercentThreshold
}

// hystrix_lowerRecallNum
func (s *flagsHystrix) setHystrixLowerRecallNum(hystrixLowerRecallNum *int) {
	s.hystrixLowerRecallNum = hystrixLowerRecallNum
}

func (s *flagsHystrix) GetHystrixLowerRecallNum() *int {
	return s.hystrixLowerRecallNum
}

// hystrix_lowerRankNum
func (s *flagsHystrix) setHystrixLowerRankNum(hystrixLowerRankNum *int) {
	s.hystrixLowerRankNum = hystrixLowerRankNum
}

func (s *flagsHystrix) GetHystrixLowerRankNum() *int {
	return s.hystrixLowerRankNum
}
