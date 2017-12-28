package countlog

import (
	"os"
	"fmt"
	"bytes"
)

var defaultLogWriter = &directLogWriter{
	logFormatter: &HumanReadableFormat{},
}

type directLogWriter struct {
	logFormatter LogFormatter
}

func (logWriter *directLogWriter) WriteLog(level int, event string, properties []interface{}) {
	if level <= LevelTrace {
		return
	}
	msg := logWriter.logFormatter.FormatLog(Event{Level: level, Event: event, Properties: properties})
	os.Stdout.Write(withColorLevelPrefix(level, msg))
	os.Stdout.Sync()
}

func withColorLevelPrefix(level int, msg []byte) []byte {
	levelColor := getColor(level)
	// ESC = \x1b
	// ESC+[ =  Control Sequence Introducer
	// \x1b[%d;1m, eg. \x1b32;1m
	// \x1b[0m
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "\x1b[%d;1m[%s]\x1b[0m%s", levelColor, getLevelName(level), msg)
	return buf.Bytes()
}

const (
	nocolor = 0
	black   = 30
	red     = 31
	green   = 32
	yellow  = 33
	blue    = 34
	purple  = 35
	cyan    = 36
	gray    = 37
)

func getColor(level int) int {
	switch level {
	case LevelTrace: return cyan
	case LevelDebug: return gray
	case LevelInfo: return green
	case LevelWarn: return yellow
	case LevelError: return red
	case LevelFatal: return purple
	default: return nocolor
	}
}