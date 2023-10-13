package flags

import "flag"

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
func (s *flagsHystrix) setHystrixTimeoutMs() {
	conf := flag.Int("hystrix_timeoutMS", 100, "")
	s.hystrixTimeoutMs = conf
}

func (s *flagsHystrix) GetHystrixTimeoutMs() *int {
	return s.hystrixTimeoutMs
}

// hystrix_MaxConcurrentRequests
func (s *flagsHystrix) setHystrixMaxConcurrentRequests() {
	conf := flag.Int("hystrix_timeoutMS", 10000, "")
	s.hystrixMaxConcurrentRequests = conf
}

func (s *flagsHystrix) GetHystrixMaxConcurrentRequests() *int {
	return s.hystrixMaxConcurrentRequests
}

// hystrix_RequestVolumeThreshold
func (s *flagsHystrix) setHystrixRequestVolumeThreshold() {
	conf := flag.Int("hystrix_RequestVolumeThreshold", 50000, "")
	s.hystrixRequestVolumeThreshold = conf
}

func (s *flagsHystrix) GetHystrixRequestVolumeThreshold() *int {
	return s.hystrixRequestVolumeThreshold
}

// hystrix_timeoutMS
func (s *flagsHystrix) setHystrixSleepWindow() {
	conf := flag.Int("hystrix_SleepWindow", 10000, "")
	s.hystrixSleepWindow = conf
}

func (s *flagsHystrix) GetHystrixSleepWindow() *int {
	return s.hystrixSleepWindow
}

// hystrix_ErrorPercentThreshold
func (s *flagsHystrix) setHystrixErrorPercentThreshold() {
	conf := flag.Int("hystrix_ErrorPercentThreshold", 5, "")
	s.hystrixTimeoutMs = conf
}

func (s *flagsHystrix) GetHystrixErrorPercentThreshold() *int {
	return s.hystrixTimeoutMs
}

// hystrix_lowerRecallNum
func (s *flagsHystrix) setHystrixLowerRecallNum() {
	conf := flag.Int("hystrix_lowerRecallNum", 100, "")
	s.hystrixLowerRecallNum = conf
}

func (s *flagsHystrix) GetHystrixLowerRecallNum() *int {
	return s.hystrixLowerRecallNum
}

// hystrix_lowerRankNum
func (s *flagsHystrix) setHystrixLowerRankNum() {
	conf := flag.Int("hystrix_lowerRankNum", 100, "")
	s.hystrixLowerRankNum = conf
}

func (s *flagsHystrix) GetHystrixLowerRankNum() *int {
	return s.hystrixLowerRankNum
}
