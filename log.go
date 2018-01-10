package plz

import (
	"github.com/v2pro/plz/countlog"
)

var LogLevel int
var LogFile string
var LogFormat countlog.LogFormatter

func setupLogging() {
	if LogLevel == 0 {
		LogLevel = countlog.LevelTrace
	}
	if LogFile == "" {
		LogFile = "STDOUT"
	}
	if LogFormat == nil {
		LogFormat = &countlog.HumanReadableFormat{}
	}
	countlog.MinLevel = LogLevel
	logWriter := countlog.NewAsyncLogWriter(
		LogLevel,
		countlog.NewFileLogOutput(LogFile))
	logWriter.LogFormat = LogFormat
	logWriter.Start()
	countlog.LogWriters = append(countlog.LogWriters, logWriter)
}
