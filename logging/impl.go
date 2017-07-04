package logging

import (
	"fmt"
	"os"
)

type placeholder struct {
	loggerKv        []interface{}
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
	got := LoggerOf(logger.loggerKv...)
	if _, stillPlaceholder := got.(*placeholder); stillPlaceholder {
		fmt.Fprintln(os.Stderr, "logger not defined yet, please add provider to logging.Providers")
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

