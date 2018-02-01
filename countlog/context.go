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
	log(LevelTrace, event, ctx, nil, properties)
}

func (ctx *Context) TraceCall(event string, err error, properties ...interface{}) {
	if err != nil {
		log(LevelError, event, ctx, err, properties)
		return
	}
	if LevelTrace < spi.MinLevel {
		return
	}
	log(LevelTrace, event, ctx, err, properties)
}

func (ctx *Context) Debug(event string, properties ...interface{}) {
	if LevelDebug < spi.MinLevel {
		return
	}
	log(LevelDebug, event, ctx, nil, properties)
}

func (ctx *Context) DebugCall(event string, err error, properties ...interface{}) {
	if err != nil {
		log(LevelError, event, ctx, err, properties)
		return
	}
	if LevelDebug < spi.MinLevel {
		return
	}
	log(LevelDebug, event, ctx, err, properties)
}

func (ctx *Context) Info(event string, properties ...interface{}) {
	if LevelInfo < spi.MinLevel {
		return
	}
	log(LevelInfo, event, ctx, nil, properties)
}

func (ctx *Context) InfoCall(event string, err error, properties ...interface{}) {
	if err != nil {
		log(LevelError, event, ctx, err, properties)
		return
	}
	if LevelInfo < spi.MinLevel {
		return
	}
	log(LevelInfo, event, ctx, err, properties)
}

func (ctx *Context) Warn(event string, properties ...interface{}) {
	log(LevelWarn, event, ctx, nil, properties)
}

func (ctx *Context) Error(event string, properties ...interface{}) {
	log(LevelError, event, ctx, nil, properties)
}

func (ctx *Context) Fatal(event string, properties ...interface{}) {
	log(LevelFatal, event, ctx, nil, properties)
}
