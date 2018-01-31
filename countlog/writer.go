package countlog

import (
	"github.com/v2pro/plz/countlog/core"
	"io"
	"sync"
)

type eventWriter struct {
	format core.Format
	writer io.Writer
}

type EventWriterConfig struct {
	Format   core.Format
	Writer   io.Writer
	Executor Executor
}

func NewEventWriter(cfg EventWriterConfig) EventSink {
	var writer io.Writer = &recylceWriter{cfg.Writer}
	if cfg.Executor != nil {
		writer = newAsyncWriter(cfg.Executor, writer)
	}
	return &eventWriter{
		format:    cfg.Format,
		writer:    writer,
	}
}

func (sink *eventWriter) HandlerOf(site *core.EventSite) core.EventHandler {
	formatter := sink.format.FormatterOf(site)
	return &writeEvent{
		formatter: formatter,
		writer:    sink.writer,
	}
}

type writeEvent struct {
	formatter core.Formatter
	writer    io.Writer
}

func (handler *writeEvent) Handle(event *core.Event) {
	space := bufPool.Get().([]byte)[:0]
	formatted := handler.formatter.Format(space, event)
	_, err := handler.writer.Write(formatted)
	if err != nil {
		// TODO: show error
	}
}

var bufPool = &sync.Pool{
	New: func() interface{} {
		return make([]byte, 0, 128)
	},
}

type recylceWriter struct {
	writer io.Writer
}

func (writer *recylceWriter) Write(buf []byte) (int, error) {
	n, err := writer.writer.Write(buf)
	bufPool.Put(buf)
	return n, err
}

//
//import (
//	"bytes"
//	"fmt"
//	"os"
//)
//
//var defaultLogWriter = &directLogWriter{
//	logFormatter: &HumanReadableFormat{},
//}
//
//type directLogWriter struct {
//	logFormatter LogFormatter
//}
//
//func (logWriter *directLogWriter) WriteLog(level int, event string, properties []interface{}) {
//	msg := logWriter.logFormatter.FormatLog(Event{Level: level, Event: event, Properties: properties})
//	if msg == nil {
//		return
//	}
//	os.Stdout.Write(withColorLevelPrefix(level, msg))
//	os.Stdout.Sync()
//}
//
//func withColorLevelPrefix(level int, msg []byte) []byte {
//	levelColor := getColor(level)
//	// ESC = \x1b
//	// ESC+[ =  Control Sequence Introducer
//	// \x1b[%d;1m, eg. \x1b32;1m
//	// \x1b[0m
//	buf := &bytes.Buffer{}
//	fmt.Fprintf(buf, "\x1b[%d;1m[%s]\x1b[0m%s", levelColor, getLevelName(level), msg)
//	return buf.Bytes()
//}
//
//const (
//	nocolor = 0
//	black   = 30
//	red     = 31
//	green   = 32
//	yellow  = 33
//	blue    = 34
//	purple  = 35
//	cyan    = 36
//	gray    = 37
//)
//
//func getColor(level int) int {
//	switch level {
//	case LevelTrace:
//		return cyan
//	case LevelDebug:
//		return gray
//	case LevelInfo:
//		return green
//	case LevelWarn:
//		return yellow
//	case LevelError:
//		return red
//	case LevelFatal:
//		return purple
//	default:
//		return nocolor
//	}
//}
