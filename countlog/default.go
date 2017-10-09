package countlog

import (
	"os"
)

var defaultLogWriter = &directLogWriter{
	logFormatter: &HumanReadableFormat{},
}

type directLogWriter struct {
	logFormatter LogFormatter
}

func (logWriter *directLogWriter) WriteLog(level int, event string, properties []interface{}) {
	msg := logWriter.logFormatter.FormatLog(Event{Level: level, Event: event, Properties: properties})
	os.Stdout.Write([]byte(msg))
	os.Stdout.Sync()
}
