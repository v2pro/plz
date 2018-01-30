package countlog

import "unsafe"

var EventSinks = []EventSink{}

type EventSink interface {
	HandlerOf(level int, eventOrCallee string,
		callerFile string, callerLine int, sample []interface{}) EventHandler
	ShouldLog(level int, eventOrCallee string,
		sample []interface{}) bool
}

type EventHandler interface {
	Handle(ctx *Context, err error, properties []interface{})
}

type EventHandlers []EventHandler

func (handlers EventHandlers) Handle(ctx *Context, err error, properties []interface{}) {
	for _, handler := range handlers {
		handler.Handle(ctx, err, properties)
	}
}

type Format interface {
	FormatterOf(level int, eventOrCallee string,
		callerFile string, callerLine int, sample []interface{}) *DummyFormatter
}

type Formatter interface {
	Format(space []byte, ctx *Context, err error, properties []interface{}) []byte
}

type Formatters []Formatter

func (formatters Formatters) Format(space []byte, ctx *Context, err error, properties []interface{}) []byte {
	for _, formatter := range formatters {
		space = formatter.Format(space, ctx, err, properties)
	}
	return space
}

type DummyFormatter struct {
	Formatter Formatter
}

func (formatter *DummyFormatter) Format(space []byte, ctx *Context, err error, properties []interface{}) []byte {
	ptr := unsafe.Pointer(&properties)
	return formatter.Formatter.Format(space, ctx, err, castEmptyInterfaces(uintptr(ptr)))
}


