package countlog

import (
	"github.com/v2pro/plz/countlog/core"
	"os"
	"github.com/v2pro/plz/countlog/output"
	"github.com/v2pro/plz/countlog/output/compact"
)

var EventSinks = []EventSink{}

type EventSink interface {
	HandlerOf(site *core.LogSite) core.EventHandler
}

// DevelopmentEventSink is used to for unit test
var DevelopmentEventSink = output.NewEventWriter(output.EventWriterConfig{
	Format: &compact.Format{},
	Writer: os.Stdout,
})


