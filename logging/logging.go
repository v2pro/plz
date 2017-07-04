package logging

import (
	"math"
)

var Providers = []func(loggerKv []interface{}) Logger{}

func LoggerOf(loggerKv ...interface{}) Logger {
	var logger Logger
	for _, provider := range Providers {
		provided := provider(loggerKv)
		if provided == nil {
			continue
		}
		logger = combineLoggers(logger, provided)
	}
	if len(Providers) == 0 {
		return &placeholder{loggerKv, nil}
	}
	if logger == nil {
		logger = &dummyLogger{}
	}
	return logger
}

type Level struct {
	Severity  int32
	LevelName string
}

var LEVEL_UNDEF = Level{math.MaxInt32, ""}
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
