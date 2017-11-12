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
	msg := logWriter.logFormatter.FormatLog(Event{Level: level, Event: event, Properties: properties})
	levelColor := getColor(level)
	// ESC = \x1b
	// ESC+[ =  Control Sequence Introducer
	// \x1b[%d;1m, eg. \x1b32;1m
	// \x1b[0m
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "\x1b[%d;1m[%s]\x1b[0m%s", levelColor, getLevelName(level), msg)
	os.Stdout.Write(buf.Bytes())
	os.Stdout.Sync()
}

const (
	nocolor = 0
	red     = 31
	green   = 32
	yellow  = 33
	blue    = 34
	gray    = 37
)

func getColor(level int) int {
	var levelColor int
	switch level {
	case LEVEL_DEBUG:
		levelColor = gray
	case LEVEL_WARN:
		levelColor = yellow
	case LEVEL_ERROR, LEVEL_FATAL:
		levelColor = red
	default:
		levelColor = green
	}
	return levelColor
}