package countlog

import (
	"github.com/v2pro/plz/countlog/core"
	"os"
	"github.com/v2pro/plz/countlog/compact"
)

var EventSinks = []EventSink{}

type EventSink interface {
	HandlerOf(level int, eventOrCallee string,
		callerFile string, callerLine int, sample []interface{}) core.EventHandler
	ShouldLog(level int, eventOrCallee string,
		sample []interface{}) bool
}

var DefaultEventSink = NewEventWriter(EventWriterConfig{
	Format: &compact.Format{},
	Writer: os.Stdout,
})


