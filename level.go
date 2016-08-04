package logging

import (
	"strings"
)

type logLevel uint8

const (
	DEBUG   logLevel = 0
	TRACE   logLevel = 1
	INFO    logLevel = 2
	WARNING logLevel = 3
	ERROR   logLevel = 4
	FATAL   logLevel = 5
	DISABLE logLevel = 255
)

func StringToLogLevel(s string) logLevel {
	switch strings.ToUpper(s) {
	case "DEBUG":
		return DEBUG
	case "TRACE":
		return TRACE
	case "INFO":
		return INFO
	case "WARN", "WARNING":
		return WARNING
	case "ERROR":
		return ERROR
	case "FATAL":
		return FATAL
	default:
		return DISABLE
	}
}

func GetLogLevelInt(s string) int {
	switch strings.ToUpper(s) {
	case "DEBUG":
		return int(DEBUG)
	case "TRACE":
		return int(TRACE)
	case "INFO":
		return int(INFO)
	case "WARN", "WARNING":
		return int(WARNING)
	case "ERROR":
		return int(ERROR)
	case "FATAL":
		return int(FATAL)
	default:
		return int(DISABLE)
	}
}

func GetLogLevelStr(level int) string {
	switch level {
	case int(DEBUG):
		return "DEBUG"
	case int(TRACE):
		return "TRACE"
	case int(INFO):
		return "INFO"
	case int(WARNING):
		return "WARNING"
	case int(ERROR):
		return "ERROR"
	case int(FATAL):
		return "FATAL"
	default:
		return ""
	}
}

func (level *logLevel) String() string {
	switch *level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO "
	case WARNING:
		return "WARN "
	case ERROR:
		return "ERROR"
	default:
		return "DISABLE"
	}
}

type levelRange struct {
	minLevel logLevel
	maxLevel logLevel
}

func (lr *levelRange) contains(level logLevel) bool {
	return level >= lr.minLevel && level <= lr.maxLevel
}
