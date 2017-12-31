package countlog

import (
	"context"
)

type Context interface {
	context.Context
	Trace(event string, properties ...interface{})
	TraceCall(callee string, err error, properties ...interface{})
	Debug(event string, properties ...interface{})
	DebugCall(callee string, err error, properties ...interface{})
	Info(event string, properties ...interface{})
	InfoCall(callee string, err error, properties ...interface{})
	Warn(event string, properties ...interface{})
	Error(event string, properties ...interface{})
	Fatal(event string, properties ...interface{})
	Log(level int, event string, properties ...interface{})
}

func Ctx(ctx context.Context) Context {
	return wrappedContext{Context: ctx}
}

type wrappedContext struct {
	context.Context
}

func (ctx wrappedContext) Trace(event string, properties ...interface{}) {
	properties = append(properties, "ctx", ctx)
	Trace(event, properties...)
}

func (ctx wrappedContext) TraceCall(callee string, err error, properties ...interface{}) {
	properties = append(properties, "ctx", ctx)
	TraceCall(callee, err, properties...)
}

func (ctx wrappedContext) Debug(event string, properties ...interface{}) {
	properties = append(properties, "ctx", ctx)
	Debug(event, properties...)
}

func (ctx wrappedContext) DebugCall(callee string, err error, properties ...interface{}) {
	properties = append(properties, "ctx", ctx)
	DebugCall(callee, err, properties...)
}

func (ctx wrappedContext) Info(event string, properties ...interface{}) {
	properties = append(properties, "ctx", ctx)
	Info(event, properties...)
}

func (ctx wrappedContext) InfoCall(callee string, err error, properties ...interface{}) {
	properties = append(properties, "ctx", ctx)
	InfoCall(callee, err, properties...)
}

func (ctx wrappedContext) Warn(event string, properties ...interface{}) {
	properties = append(properties, "ctx", ctx)
	Warn(event, properties...)
}

func (ctx wrappedContext) Error(event string, properties ...interface{}) {
	properties = append(properties, "ctx", ctx)
	Error(event, properties...)
}

func (ctx wrappedContext) Fatal(event string, properties ...interface{}) {
	properties = append(properties, "ctx", ctx)
	Fatal(event, properties...)
}

func (ctx wrappedContext) Log(level int, event string, properties ...interface{}) {
	properties = append(properties, "ctx", ctx)
	Log(level, event, properties...)
}
