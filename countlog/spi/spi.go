package spi

import (
	"context"
	"time"
)

// MinLevel exists to minimize the overhead of Trace/Debug logging
var MinLevel = LevelTrace
var SuccinctLevel = LevelDebug

const LevelTrace = 10
const LevelDebug = 20
const LevelInfo = 30
const LevelWarn = 40
const LevelError = 50
const LevelFatal = 60

// LogSite is the location of log in the source code
type LogSite struct {
	File   string
	Line   int
	// Level might change, the actual event can have different level
	Level  int
	Event  string
	Agg    string
	Sample []interface{}
}

type Event struct {
	Level      int
	Context    context.Context
	Error      error
	Timestamp  time.Time
	Properties []interface{}
}

type EventSink interface {
	HandlerOf(site *LogSite) EventHandler
}

type EventHandler interface {
	Handle(event *Event)
}

type EventHandlers []EventHandler

func (handlers EventHandlers) Handle(event *Event) {
	for _, handler := range handlers {
		handler.Handle(event)
	}
}

type SelectiveEventHandler struct {
	Verbose  EventHandler
	Succinct EventHandler
}

func (handler *SelectiveEventHandler) Handle(event *Event) {
	if event.Level >= SuccinctLevel {
		handler.Succinct.Handle(event)
	}
	handler.Verbose.Handle(event)
}

type SelectiveEventSink struct {
	Verbose  EventSink
	Succinct EventSink
}

func (sink *SelectiveEventSink) HandlerOf(site *LogSite) EventHandler {
	verbose := sink.Verbose.HandlerOf(site)
	if verbose == nil {
		return sink.Succinct.HandlerOf(site)
	}
	return &SelectiveEventHandler{
		Verbose:  verbose,
		Succinct: sink.Succinct.HandlerOf(site),
	}
}
