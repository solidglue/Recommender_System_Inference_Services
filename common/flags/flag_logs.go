package flags

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
func (s *flagsLog) setLogMaxSize(logMaxSize *int) {
	s.logMaxSize = logMaxSize
}

func (s *flagsLog) GetLogMaxSize() *int {
	return s.logMaxSize
}

// log_save_days
func (s *flagsLog) setLogSaveDays(logSaveDays *int) {
	s.logSaveDays = logSaveDays
}

func (s *flagsLog) GetLogSaveDays() *int {
	return s.logSaveDays
}

// log_file_name
func (s *flagsLog) setLogFileName(logFileName *string) {
	s.logFileName = logFileName
}

func (s *flagsLog) GetLogFileName() *string {
	return s.logFileName
}

// log_level
func (s *flagsLog) setLogLevel(logLevel *string) {
	s.logLevel = logLevel
}

func (s *flagsLog) GetLogLevel() *string {
	return s.logLevel
}
