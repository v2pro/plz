package log

import (
	"math"
	"fmt"
	"os"
)

var providers = []func(loggerKv []interface{}) Logger{}

var GetLogger = func(loggerKv ...interface{}) Logger {
	var logger Logger
	for _, provider := range providers {
		provided := provider(loggerKv)
		if provided == nil {
			continue
		}
		logger = combineLoggers(logger, provided)
	}
	if len(providers) == 0 {
		return &placeholder{loggerKv, nil}
	}
	if logger == nil {
		logger = &dummyLogger{}
	}
	return logger
}

func AddLoggerProvider(provider func(loggerKv []interface{}) Logger) {
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
	Log(level Level, msg string, kv ...interface{})
	Error(msg string, kv ...interface{})
	Info(msg string, kv ...interface{})
	Debug(msg string, kv ...interface{})
	ShouldLog(level Level) bool
}

type placeholder struct{
	loggerKv []interface{}
	realLoggerCache Logger
}

func (logger *placeholder) Log(level Level, msg string, kv ...interface{}) {
	logger.realLogger().Log(level, msg, kv...)
}
func (logger *placeholder) Error(msg string, kv ...interface{}) {
	logger.realLogger().Error(msg, kv...)
}
func (logger *placeholder) Info(msg string, kv ...interface{}) {
	logger.realLogger().Info(msg, kv...)
}
func (logger *placeholder) Debug(msg string, kv ...interface{}) {
	logger.realLogger().Debug(msg, kv...)
}
func (logger *placeholder) ShouldLog(level Level) bool {
	return logger.realLogger().ShouldLog(level)
}
func (logger *placeholder) realLogger() Logger {
	if logger.realLoggerCache != nil {
		return logger.realLoggerCache
	}
	got := GetLogger(logger.loggerKv...)
	if _, stillPlaceholder := got.(*placeholder); stillPlaceholder {
		fmt.Fprintln(os.Stderr, "logger not defined yet, please AddLoggerProvider")
		return &dummyLogger{}
	}
	logger.realLoggerCache = got
	return got
}

type dummyLogger struct{}

func (logger *dummyLogger) Log(level Level, msg string, kv ...interface{}) {
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

func (logger *combinedLogger) Log(level Level, msg string, kv ...interface{}) {
	for _, logger := range logger.loggers {
		logger.Log(level, msg, kv...)
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
	loggerKv string
	minLevel Level
}

func (logger *stderrLogger) Log(level Level, msg string, kv ...interface{}) {
	fmt.Fprintln(os.Stderr, append([]interface{}{logger.loggerKv, level.LevelName, msg}, kv...)...)
}
func (logger *stderrLogger) Error(msg string, kv ...interface{}) {
	logger.Log(LEVEL_ERROR, msg, kv...)
}
func (logger *stderrLogger) Info(msg string, kv ...interface{}) {
	logger.Log(LEVEL_INFO, msg, kv...)
}
func (logger *stderrLogger) Debug(msg string, kv ...interface{}) {
	logger.Log(LEVEL_DEBUG, msg, kv...)
}
func (logger *stderrLogger) ShouldLog(level Level) bool {
	return level.Severity >= logger.minLevel.Severity
}

func NewStderrLogger(loggerKv []interface{}, minLevel Level) Logger {
	return &stderrLogger{fmt.Sprintf("%v", loggerKv), minLevel}
}
