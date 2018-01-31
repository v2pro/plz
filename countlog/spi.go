package countlog

import (
	"github.com/v2pro/plz/countlog/core"
	"os"
	"github.com/v2pro/plz/countlog/compact"
)

var EventSinks = []EventSink{}

type EventSink interface {
	HandlerOf(site *core.EventSite) core.EventHandler
}

var DefaultEventSink = NewEventWriter(EventWriterConfig{
	Format: &compact.Format{},
	Writer: os.Stdout,
})


