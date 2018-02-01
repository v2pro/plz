package spi

import (
	"context"
	"time"
)

// MinLevel exists to minimize the overhead of Trace/Debug logging
var MinLevel = LevelTrace
// Succinct event handler need level above SuccinctLevel to output
var SuccinctLevel = LevelDebugCall

const LevelTraceCall = 5

// this should be development default
const LevelTrace = 10
const LevelDebugCall = 15
const LevelDebug = 20
const LevelInfoCall = 25

// this should be the production default
const LevelInfo = 30

// LevelWarn is the level for error != nil
const LevelWarn = 40

// LevelError is the level for user visible error
const LevelError = 50

// LevelFatal is the level for panic or panic like scenario
const LevelFatal = 60

func LevelName(level int) string {
	switch level {
	case LevelTraceCall:
		return "TRACE_CALL"
	case LevelTrace:
		return "TRACE"
	case LevelDebugCall:
		return "DEBUG_CALL"
	case LevelDebug:
		return "DEBUG"
	case LevelInfoCall:
		return "INFO_CALL"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// LogSite is the location of log in the source code
type LogSite struct {
	File string
	Line int
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

type DummyEventHandler struct {
}

func (handler *DummyEventHandler) Handle(event *Event) {
}
