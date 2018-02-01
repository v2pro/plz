package countlog

import (
	"context"
	"github.com/v2pro/plz/countlog/spi"
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

func (ctx *Context) Trace(event string, properties ...interface{}) {
	if LevelTrace < spi.MinLevel {
		return
	}
	log(LevelTrace, event, "", ctx, nil, properties)
}

func (ctx *Context) TraceCall(event string, err error, properties ...interface{}) error {
	if err != nil {
		return log(LevelWarn, event, "call", ctx, err, properties)
	}
	if LevelTrace < spi.MinLevel {
		return nil
	}
	log(LevelTrace, event, "call", ctx, err, properties)
	return nil
}

func (ctx *Context) Debug(event string, properties ...interface{}) {
	if LevelDebug < spi.MinLevel {
		return
	}
	log(LevelDebug, event, "", ctx, nil, properties)
}

func (ctx *Context) DebugCall(event string, err error, properties ...interface{}) error {
	if err != nil {
		return log(LevelWarn, event, "call", ctx, err, properties)
	}
	if LevelDebug < spi.MinLevel {
		return nil
	}
	log(LevelDebug, event, "call", ctx, err, properties)
	return nil
}

func (ctx *Context) Info(event string, properties ...interface{}) {
	if LevelInfo < spi.MinLevel {
		return
	}
	log(LevelInfo, event, "", ctx, nil, properties)
}

func (ctx *Context) InfoCall(event string, err error, properties ...interface{}) error {
	if err != nil {
		return log(LevelWarn, event, "call", ctx, err, properties)
	}
	if LevelInfo < spi.MinLevel {
		return nil
	}
	log(LevelInfo, event, "call", ctx, err, properties)
	return nil
}

func (ctx *Context) Warn(event string, properties ...interface{}) {
	log(LevelWarn, event, "", ctx, nil, properties)
}

func (ctx *Context) Error(event string, properties ...interface{}) {
	log(LevelError, event, "", ctx, nil, properties)
}

func (ctx *Context) Fatal(event string, properties ...interface{}) {
	log(LevelFatal, event, "", ctx, nil, properties)
}
