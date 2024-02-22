package logs

import (
	"fmt"
	"runtime"

	//"runtime"
	"strconv"
	"strings"
	"sync"
)

const asyncLogBuffer = 1000

var once sync.Once

func InitLogParam(file_name string, log_level string) error {
	once.Do(func() {
		// var maxSize = *config.FLAG_log_max_size
		// var maxDays = *config.FLAG_log_keep_days

		var maxSize = 10
		var maxDays = 7

		var filename = file_name
		var logLevelStr = log_level
		logLevel := LevelInfo
		if strings.EqualFold(logLevelStr, "fatal") {
			logLevel = LevelFatal
		} else if strings.EqualFold(logLevelStr, "error") {
			logLevel = LevelError
		} else if strings.EqualFold(logLevelStr, "debug") {
			logLevel = LevelDebug
		}

		lb := GetBeeLogger()

		fileLogConfig := fmt.Sprintf(`{"filename":"%v","maxsize":%v,
                "daily":true,"maxdays":%v, "rotate":true,
                "level":%v,"perm": "0666","separate":["error", "warning", "info", "debug"]}`,
			filename, maxSize, maxDays, LevelTrace)
		// AdapterMultiFile
		if err := lb.SetLogger(AdapterMultiFile, fileLogConfig); err != nil {
			panic("Logs module Init Failed")
		}

		lb.SetLevel(logLevel)
		lb.SetLogFuncCallDepth(3)
		lb.SetShortFile(true)
		lb.Async(asyncLogBuffer)
	})
	return nil
}

func InitLog() error {
	once.Do(func() {
		// var maxSize = *config.FLAG_log_max_size
		// var maxDays = *config.FLAG_log_keep_days
		// var filename = *config.FLAG_log_file_name
		// var logLevelStr = *config.FLAG_log_level
		// logLevel := LevelInfo
		// if strings.EqualFold(logLevelStr, "fatal") {
		// 	logLevel = LevelFatal
		// } else if strings.EqualFold(logLevelStr, "error") {
		// 	logLevel = LevelError
		// } else if strings.EqualFold(logLevelStr, "debug") {
		// 	logLevel = LevelDebug
		// }

		// lb := GetBeeLogger()

		// fileLogConfig := fmt.Sprintf(`{"filename":"%v","maxsize":%v,
		//         "daily":true,"maxdays":%v, "rotate":true,
		//         "level":%v,"perm": "0666","separate":["error", "warning", "info", "debug"]}`,
		// 	filename, maxSize, maxDays, LevelTrace)
		// // AdapterMultiFile
		// if err := lb.SetLogger(AdapterMultiFile, fileLogConfig); err != nil {
		// 	panic("Logs module Init Failed")
		// }

		// lb.SetLevel(logLevel)
		// lb.SetLogFuncCallDepth(3)
		// lb.SetShortFile(true)
		// lb.Async(asyncLogBuffer)
	})
	return nil
}

func CloseLog() {
	GetBeeLogger().Close()
	beeLogger = NewLogger()
}

func Close() {
	beeLogger.Close()
	beeLogger = NewLogger()
}

// SetLogger sets a new logger.
func SetLogger(adapter string, config ...string) error {
	return beeLogger.SetLogger(adapter, config...)
}

// Emergency logs a message at emergency level.
func Emergency(f interface{}, v ...interface{}) {
	beeLogger.Emergency(formatLog(f, v...))
}

// Alert logs a message at alert level.
func Alert(f interface{}, v ...interface{}) {
	beeLogger.Alert(formatLog(f, v...))
}

// Critical logs a message at critical level.
func Critical(f interface{}, v ...interface{}) {
	beeLogger.Critical(formatLog(f, v...))
}

// Fatal logs a message at fatal level.
func Fatal(f interface{}, v ...interface{}) {
	log := formatLog(f, v...)
	beeLogger.Fatal(log)

	_, filename, line, ok := runtime.Caller(beeLogger.loggerFuncCallDepth - 1)
	if !ok {
		filename = "???"
		line = 0
	} else if beeLogger.shortFile {
		for i := len(filename) - 1; i >= 0; i-- {
			if filename[i] == '/' {
				filename = filename[i+1:]
				break
			}
		}
	}
	msg := filename + ":" + strconv.Itoa(line)
	log = log + "|" + msg
}

// Error logs a message at error level.
func Error(f interface{}, v ...interface{}) {
	log := formatLog(f, v...)
	beeLogger.Error(log)

	_, filename, line, ok := runtime.Caller(beeLogger.loggerFuncCallDepth - 1)
	if !ok {
		filename = "???"
		line = 0
	} else if beeLogger.shortFile {
		for i := len(filename) - 1; i >= 0; i-- {
			if filename[i] == '/' {
				filename = filename[i+1:]
				break
			}
		}
	}
	msg := filename + ":" + strconv.Itoa(line)
	log = log + "|" + msg
}

// Warning logs a message at warning level.
func Warning(f interface{}, v ...interface{}) {
	beeLogger.Warn(formatLog(f, v...))
}

// Warn compatibility alias for Warning()
func Warn(f interface{}, v ...interface{}) {
	beeLogger.Warn(formatLog(f, v...))
}

// Notice logs a message at notice level.
func Notice(f interface{}, v ...interface{}) {
	beeLogger.Notice(formatLog(f, v...))
}

// Informational logs a message at info level.
func Informational(f interface{}, v ...interface{}) {
	beeLogger.Informational(formatLog(f, v...))
}

// Info compatibility alias for Warning()
func Info(f interface{}, v ...interface{}) {
	beeLogger.Info(formatLog(f, v...))
}

// Debug logs a message at debug level.
func Debug(f interface{}, v ...interface{}) {
	beeLogger.Debug(formatLog(f, v...))
}

// Trace logs a message at trace level.
// compatibility alias for Warning()
func Trace(f interface{}, v ...interface{}) {
	beeLogger.Trace(formatLog(f, v...))
}
