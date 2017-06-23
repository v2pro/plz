package log

import (
	"math"
	"fmt"
	"os"
	"strings"
)

var providers = []func(path []string) Logger{}

var GetLogger = func(path ...string) Logger {
	var logger Logger
	for _, provider := range providers {
		provided := provider(path)
		if provided == nil {
			continue
		}
		logger = combineLoggers(logger, provided)
	}
	if len(providers) == 0 {
		panic("logger not defined yet, please AddLoggerProvider")
	}
	if logger == nil {
		logger = &dummyLogger{}
	}
	return logger
}

func AddLoggerProvider(provider func(path []string) Logger) {
	providers = append(providers, provider)
}

func PathStartedWith(path []string, prefix ...string) bool {
	if len(prefix) > len(path) {
		return false
	}
	for i, p := range prefix {
		if path[i] != p {
			return false
		}
	}
	return true
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
	Log(kv ...interface{})
	LogMessage(level Level, msg string, kv ...interface{})
	Error(msg string, kv ...interface{})
	Info(msg string, kv ...interface{})
	Debug(msg string, kv ...interface{})
	ShouldLog(level Level) bool
}

type dummyLogger struct{}

func (logger *dummyLogger) Log(kv ...interface{}) {
}
func (logger *dummyLogger) LogMessage(level Level, msg string, kv ...interface{}) {
}
func (logger *dummyLogger) Error(msg string, kv ...interface{}) {
}
func (logger *dummyLogger) Info(msg string, kv ...interface{}) {
}
func (logger *dummyLogger) Debug(msg string, kv ...interface{}) {
}
func (logger *dummyLogger) ShouldLog(level Level) bool {
	return false
}

type combinedLogger struct {
	loggers []Logger
}

func (logger *combinedLogger) Log(kv ...interface{}) {
	for _, logger := range logger.loggers {
		logger.Log(kv...)
	}
}
func (logger *combinedLogger) LogMessage(level Level, msg string, kv ...interface{}) {
	for _, logger := range logger.loggers {
		logger.LogMessage(level, msg, kv...)
	}
}
func (logger *combinedLogger) Error(msg string, kv ...interface{}) {
	for _, logger := range logger.loggers {
		logger.Error(msg, kv...)
	}
}
func (logger *combinedLogger) Info(msg string, kv ...interface{}) {
	for _, logger := range logger.loggers {
		logger.Info(msg, kv...)
	}
}
func (logger *combinedLogger) Debug(msg string, kv ...interface{}) {
	for _, logger := range logger.loggers {
		logger.Debug(msg, kv...)
	}
}
func (logger *combinedLogger) ShouldLog(level Level) bool {
	for _, logger := range logger.loggers {
		if logger.ShouldLog(level) {
			return true
		}
	}
	return false
}
func combineLoggers(oldLogger Logger, newLogger Logger) Logger {
	if oldLogger == nil {
		return newLogger
	}
	asCombinedLogger, ok := oldLogger.(*combinedLogger)
	if ok {
		asCombinedLogger.loggers = append(asCombinedLogger.loggers, newLogger)
		return asCombinedLogger
	}
	return &combinedLogger{loggers: []Logger{oldLogger, newLogger}}
}

type stderrLogger struct {
	loggerName string
	minLevel   Level
}

func (logger *stderrLogger) Log(kv ...interface{}) {
	logger.LogMessage(LEVEL_UNDEF, "", kv...)
}
func (logger *stderrLogger) LogMessage(level Level, msg string, kv ...interface{}) {
	fmt.Fprintln(os.Stderr, append([]interface{}{logger.loggerName, level.LevelName, msg}, kv...)...)
}
func (logger *stderrLogger) Error(msg string, kv ...interface{}) {
	logger.LogMessage(LEVEL_ERROR, msg, kv...)
}
func (logger *stderrLogger) Info(msg string, kv ...interface{}) {
	logger.LogMessage(LEVEL_INFO, msg, kv...)
}
func (logger *stderrLogger) Debug(msg string, kv ...interface{}) {
	logger.LogMessage(LEVEL_DEBUG, msg, kv...)
}
func (logger *stderrLogger) ShouldLog(level Level) bool {
	return level.Severity >= logger.minLevel.Severity
}

func NewStderrLogger(path []string, minLevel Level) Logger {
	return &stderrLogger{strings.Join(path, "."), minLevel}
}
