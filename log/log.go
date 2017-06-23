package log

import "math"

type Level struct {
	Severity  int32
	LevelName string
}

var LEVEL_UNDEF = Level{math.MaxInt32, "UNDEF"}
var LEVEL_FATAL = Level{60, "FATAL"}
var LEVEL_ERROR = Level{50, "ERROR"}
var LEVEL_WARNING = Level{40, "WARNING"}
var LEVEL_INFO = Level{30, "INFO"}
var LEVEL_DEBUG = Level{20, "DEBUG"}
var LEVEL_TRACE = Level{10, "TRACE"}

type Logger interface {
	Log(level Level, msg string, kv ...interface{})
	Error(msg string, kv ...interface{})
	Info(msg string, kv ...interface{})
	Debug(msg string, kv ...interface{})
	ShouldLog(level Level) bool
}