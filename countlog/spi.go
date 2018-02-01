package countlog

import (
	"github.com/v2pro/plz/countlog/core"
	"os"
	"github.com/v2pro/plz/countlog/compact"
)

var EventSinks = []EventSink{}

type EventSink interface {
	HandlerOf(site *core.LogSite) core.EventHandler
}

// DevelopmentEventSink is used to for unit test
var DevelopmentEventSink = NewEventWriter(EventWriterConfig{
	Format: &compact.Format{},
	Writer: os.Stdout,
})


