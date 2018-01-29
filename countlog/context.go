package countlog

import (
	"context"
)

func Ctx(ctx context.Context) *Context {
	wrapped, isWrapped := ctx.(*Context)
	if isWrapped {
		return wrapped
	}
	return &Context{Context: ctx}
}

type Context struct {
	context.Context
}

func (ctx *Context) TraceCall(callee string, err error, properties ...interface{}) {
	properties = append(properties, "ctx", ctx)
	TraceCall(callee, err, properties...)
}

func (ctx *Context) Debug(event string, properties ...interface{}) {
	properties = append(properties, "ctx", ctx)
	Debug(event, properties...)
}

func (ctx *Context) DebugCall(callee string, err error, properties ...interface{}) {
	properties = append(properties, "ctx", ctx)
	DebugCall(callee, err, properties...)
}

func (ctx Context) Info(event string, properties ...interface{}) {
	properties = append(properties, "ctx", ctx)
	Info(event, properties...)
}

func (ctx *Context) InfoCall(callee string, err error, properties ...interface{}) {
	properties = append(properties, "ctx", ctx)
	InfoCall(callee, err, properties...)
}

func (ctx *Context) Warn(event string, properties ...interface{}) {
	properties = append(properties, "ctx", ctx)
	Warn(event, properties...)
}

func (ctx *Context) Error(event string, properties ...interface{}) {
	properties = append(properties, "ctx", ctx)
	Error(event, properties...)
}

func (ctx *Context) Fatal(event string, properties ...interface{}) {
	properties = append(properties, "ctx", ctx)
	Fatal(event, properties...)
}

func (ctx *Context) Log(level int, event string, properties ...interface{}) {
	properties = append(properties, "ctx", ctx)
	Log(level, event, properties...)
}
