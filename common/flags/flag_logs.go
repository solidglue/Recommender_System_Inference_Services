package flags

import "flag"

var flagsLogInstance *flagsLog

type flagsLog struct {
	// logs
	logMaxSize  *int
	logSaveDays *int
	logFileName *string
	logLevel    *string
}

// singleton instance
func init() {
	flagsLogInstance = new(flagsLog)
}

func getFlagLogInstance() *flagsLog {
	return flagsLogInstance
}

// log_max_size
func (s *flagsLog) setLogMaxSize() {
	conf := flag.Int("log_max_size", 200000000, "the max size of the log file (in Byte)")
	s.logMaxSize = conf
}

func (s *flagsLog) GetLogMaxSize() *int {
	return s.logMaxSize
}

// log_save_days
func (s *flagsLog) setLogSaveDays() {
	conf := flag.Int("log_save_days", 7, "")
	s.logSaveDays = conf
}

func (s *flagsLog) GetLogSaveDays() *int {
	return s.logSaveDays
}

// log_file_name
func (s *flagsLog) setLogFileName() {
	conf := flag.String("log_file_name", "infer.log", "")
	s.logFileName = conf
}

func (s *flagsLog) GetLogFileName() *string {
	return s.logFileName
}

// log_level
func (s *flagsLog) setLogLevel() {
	conf := flag.String("log_level", "error", "the log level, (debug, info, error, fatal)")
	s.logLevel = conf
}

func (s *flagsLog) GetLogLevel() *string {
	return s.logLevel
}
