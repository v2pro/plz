package countlog

import (
	"context"
)

type Context interface {
	context.Context
	Trace(event string, properties ...interface{})
	TraceMetric(callee string, err error, properties ...interface{})
	Debug(event string, properties ...interface{})
	DebugMetric(callee string, err error, properties ...interface{})
	Info(event string, properties ...interface{})
	InfoMetric(callee string, err error, properties ...interface{})
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
	properties = append(properties, "ctx", ctx.Context)
	Trace(event, properties...)
}

func (ctx wrappedContext) TraceMetric(callee string, err error, properties ...interface{}) {
	properties = append(properties, "ctx", ctx.Context)
	TraceMetric(callee, err, properties...)
}

func (ctx wrappedContext) Debug(event string, properties ...interface{}) {
	properties = append(properties, "ctx", ctx.Context)
	Debug(event, properties...)
}

func (ctx wrappedContext) DebugMetric(callee string, err error, properties ...interface{}) {
	properties = append(properties, "ctx", ctx.Context)
	DebugMetric(callee, err, properties...)
}

func (ctx wrappedContext) Info(event string, properties ...interface{}) {
	properties = append(properties, "ctx", ctx.Context)
	Info(event, properties...)
}

func (ctx wrappedContext) InfoMetric(callee string, err error, properties ...interface{}) {
	properties = append(properties, "ctx", ctx.Context)
	InfoMetric(callee, err, properties...)
}

func (ctx wrappedContext) Warn(event string, properties ...interface{}) {
	properties = append(properties, "ctx", ctx.Context)
	Warn(event, properties...)
}

func (ctx wrappedContext) Error(event string, properties ...interface{}) {
	properties = append(properties, "ctx", ctx.Context)
	Error(event, properties...)
}

func (ctx wrappedContext) Fatal(event string, properties ...interface{}) {
	properties = append(properties, "ctx", ctx.Context)
	Fatal(event, properties...)
}

func (ctx wrappedContext) Log(level int, event string, properties ...interface{}) {
	properties = append(properties, "ctx", ctx.Context)
	Log(level, event, properties...)
}
