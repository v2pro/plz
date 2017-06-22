package plz

import "math"

type Logger interface {
	Log(level LogLevel, msg string, kv ...interface{})
	Error(msg string, kv ...interface{})
	Info(msg string, kv ...interface{})
	Debug(msg string, kv ...interface{})
	ShouldLog(level LogLevel) bool
}

var GetLogger func(name string) Logger

type LogLevel struct {
	Severity  int32
	LevelName string
}

var LOG_LEVEL_UNDEF = LogLevel{math.MaxInt32, "UNDEF"}
var LOG_LEVEL_FATAL = LogLevel{60, "FATAL"}
var LOG_LEVEL_ERROR = LogLevel{50, "ERROR"}
var LOG_LEVEL_WARNING = LogLevel{40, "WARNING"}
var LOG_LEVEL_INFO = LogLevel{30, "INFO"}
var LOG_LEVEL_DEBUG = LogLevel{20, "DEBUG"}
var LOG_LEVEL_TRACE = LogLevel{10, "TRACE"}
