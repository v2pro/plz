package plz

import (
	"context"
	"time"
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

func (r *Routine) goLongRunning(ctx context.Context) {
	go func() {
		for restartedTimes := 0; r.goLongRunningOnce(ctx); restartedTimes++ {
			shouldRestartAgain := OnGoroutineRestarted(r, restartedTimes)
			if !shouldRestartAgain {
				break
			}
		}
	}()
}

func (r *Routine) goLongRunningOnce(ctx context.Context) (notDone bool) {
	defer func() {
		recovered := recover()
		if recovered != nil {
			for _, handlePanic := range panicHandlers {
				handlePanic(r, recovered)
			}
			notDone = true
		}
	}()
	return r.LongRunning(ctx)
}

func (r *Routine) goOneOff() {
	go func() {
		defer func() {
			recovered := recover()
			for _, handlePanic := range panicHandlers {
				handlePanic(r, recovered)
			}
		}()
		r.OneOff()
	}()
}

func Go(oneOff func()) {
	Routine{OneOff: oneOff}.Go()
}

type HandlePanic func(routine *Routine, recovered interface{})

var panicHandlers = []HandlePanic{}
var OnGoroutineRestarted = func(routine *Routine, restartedTimes int) bool {
	time.Sleep(100 * time.Microsecond)
	return true
}

func RegisterPanicHandler(handlePanic HandlePanic) {
	panicHandlers = append(panicHandlers, handlePanic)
}
