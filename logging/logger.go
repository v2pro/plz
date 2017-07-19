package logging

import (
	"fmt"
	"os"
	"runtime"
)

type placeholder struct {
	loggerKV        []interface{}
	realLoggerCache Logger
}

func (logger *placeholder) Log(level Level, msg string, kv ...interface{}) {
	logger.realLogger().Log(level, msg, kv...)
}
func (logger *placeholder) Error(err error, msg string, kv ...interface{}) error {
	return logger.realLogger().Error(err, msg, kv...)
}
func (logger *placeholder) Warning(msg string, kv ...interface{}) {
	logger.realLogger().Warning(msg, kv...)
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
func (logger *placeholder) SetLevel(level Level) Logger {
	return logger.realLogger().SetLevel(level)
}
func (logger *placeholder) realLogger() Logger {
	if logger.realLoggerCache != nil {
		return logger.realLoggerCache
	}
	logger.realLoggerCache = realLoggerOf(logger.loggerKV)
	return logger.realLoggerCache
}

type defaultLogger struct {
	loggerKV        []interface{}
	logWriter LogWriter
	minLevel Level
}

func (logger *defaultLogger) Log(level Level, msg string, kv ...interface{}) {
	if logger.ShouldLog(level) {
		logger.logWriter.Log(level, msg, kv)
	}
}
func (logger *defaultLogger) Error(err error, msg string, kv ...interface{}) error {
	_, file, line, _ := runtime.Caller(0)
	logger.Log(ErrorLevel, msg, append(kv, "file", file, "line", line, "err", err)...)
	if err == nil {
		return nil
	}
	return &errorWrapper{cause:err, msg: msg, kv: kv}
}

func (logger *defaultLogger) Warning(msg string, kv ...interface{}) {
	logger.Log(WarningLevel, msg, kv...)
}

func (logger *defaultLogger) Info(msg string, kv ...interface{}) {
	logger.Log(InfoLevel, msg, kv...)
}
func (logger *defaultLogger) Debug(msg string, kv ...interface{}) {
	logger.Log(DebugLevel, msg, kv...)
}
func (logger *defaultLogger) ShouldLog(level Level) bool {
	return level.Severity >= logger.minLevel.Severity
}
func (logger *defaultLogger) SetLevel(level Level) Logger {
	logger.minLevel = level
	return logger
}

type combinedLogWriter struct {
	logWriters []LogWriter
}

func (logger *combinedLogWriter) Log(level Level, msg string, kv []interface{}) {
	for _, logWriter := range logger.logWriters {
		logWriter.Log(level, msg, kv)
	}
}

type stderrLogWriter struct {
}

func (logger *stderrLogWriter) Log(level Level, msg string, kv []interface{}) {
	line := []byte(level.LevelName)
	line = append(line, ' ')
	line = append(line, msg...)
	for i := 0; i < len(kv); i+=2 {
		line = append(line, fmt.Sprintf("||%v=%v", kv[i], kv[i+1])...)
	}
	fmt.Fprintln(os.Stderr, string(line))
	os.Stderr.Sync()
}
