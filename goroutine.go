package plz

import (
	"context"
)

type Routine struct {
	ParentContext context.Context
	OneOff        func()
	LongRunning   func(ctx context.Context) bool
}

func (cfg Routine) Go() context.CancelFunc {
	ctx := cfg.ParentContext
	cancel := func() {}
	if ctx != nil {
		ctx, cancel = context.WithCancel(ctx)
	}
	if cfg.OneOff != nil {
		cfg.goOneOff()
	} else {
		cfg.goLongRunning(ctx)
	}
	return cancel
}
func (cfg *Routine) goLongRunning(ctx context.Context) {
	go func() {
		for cfg.goLongRunningOnce(ctx) {
		}
	}()
}

func (cfg *Routine) goLongRunningOnce(ctx context.Context) (notDone bool) {
	defer func() {
		recovered := recover()
		if recovered != nil {
			for _, handlePanic := range panicHandlers {
				handlePanic(cfg, recovered)
			}
			notDone = true
		}
	}()
	return cfg.LongRunning(ctx)
}

func (cfg *Routine) goOneOff() {
	go func() {
		defer func() {
			recovered := recover()
			for _, handlePanic := range panicHandlers {
				handlePanic(cfg, recovered)
			}
		}()
		cfg.OneOff()
	}()
}

func Go(oneOff func()) {
	Routine{OneOff: oneOff}.Go()
}

type HandlePanic func(routine *Routine, recovered interface{})

var panicHandlers = []HandlePanic{}

func RegisterPanicHandler(handlePanic HandlePanic) {
	panicHandlers = append(panicHandlers, handlePanic)
}
