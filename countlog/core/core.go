package core

import (
	"context"
	"time"
)

// LogSite is the location of log in the source code
type LogSite struct {
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