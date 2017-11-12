package countlog

import (
	"os"
	"runtime"
	"bytes"
	"fmt"
)

type AsyncLogWriter struct {
	MinLevel       int
	EventWhitelist map[string]bool
	msgChan        chan Event
	isClosed       chan bool
	LogFormatter   LogFormatter
	LogOutput      LogOutput
}

func (logWriter *AsyncLogWriter) ShouldLog(level int, event string, properties []interface{}) bool {
	if logWriter.EventWhitelist[event] {
		return true
	}
	return level >= logWriter.MinLevel
}

func (logWriter *AsyncLogWriter) WriteLog(level int, event string, properties []interface{}) {
	select {
	case logWriter.msgChan <- Event{Level: level, Event: event, Properties: properties}:
	default:
		// drop on the floor
	}
}

func (logWriter *AsyncLogWriter) Close() {
	close(logWriter.isClosed)
	if logWriter.LogOutput != nil {
		logWriter.LogOutput.Close()
	}
}

func (logWriter *AsyncLogWriter) Start() {
	go func() {
		defer func() {
			recovered := recover()
			if recovered != nil {
				os.Stderr.WriteString(fmt.Sprintf("countlog panic: %v\n", recovered))
				buf := make([]byte, 1<<16)
				runtime.Stack(buf, true)
				os.Stderr.Write(buf)
				os.Stderr.Sync()
			}
		}()
		for {
			select {
			case event := <-logWriter.msgChan:
				formattedEvent := logWriter.LogFormatter.FormatLog(event)
				if _, ok := logWriter.LogOutput.(*osFileLogOutput); ok {
					levelColor := getColor(event.Level)
					// ESC = \x1b
					// ESC+[ =  Control Sequence Introducer
					// \x1b[%d;1m, eg. \x1b32;1m
					// \x1b[0m
					buf := &bytes.Buffer{}
					fmt.Fprintf(buf, "\x1b[%d;1m[%s]\x1b[0m%s", levelColor, getLevelName(event.Level), formattedEvent)
					logWriter.LogOutput.OutputLog(event.Properties[1].(int64), buf.Bytes())
				} else {
					logWriter.LogOutput.OutputLog(event.Properties[1].(int64), formattedEvent)
				}
			case <-logWriter.isClosed:
				return
			}
		}
	}()
}

func NewAsyncLogWriter(minLevel int, output LogOutput) *AsyncLogWriter {
	writer := &AsyncLogWriter{
		MinLevel:       minLevel,
		msgChan:        make(chan Event, 1024),
		LogFormatter:   &HumanReadableFormat{},
		LogOutput:      output,
		EventWhitelist: map[string]bool{},
	}
	return writer
}
