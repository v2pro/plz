package core

import (
	"context"
	"time"
)

// EventSite is the location of log in the source code
type EventSite struct {
	File string
	Line int
	Level int
	EventOrCallee string
	Agg string
	Sample []interface{}
}

type Event struct {
	Context context.Context
	Error error
	Timestamp time.Time
	Properties []interface{}
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

type Format interface {
	FormatterOf(site *EventSite) Formatter
}

type Formatter interface {
	Format(space []byte, event *Event) []byte
}

type Formatters []Formatter

func (formatters Formatters) Format(space []byte, event *Event) []byte {
	for _, formatter := range formatters {
		space = formatter.Format(space, event)
	}
	return space
}