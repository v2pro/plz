package routine

import (
	"context"
	"time"
)

func Go(oneOff func()) error {
	_, err := Of{OneOff: oneOff}.Go()
	return err
}

func GoLongRunning(longRunning func(ctx context.Context)) (context.CancelFunc, error) {
	return Of{LongRunning: longRunning}.Go()
}

type Of struct {
	ParentContext context.Context
	OneOff        func()
	LongRunning   func(ctx context.Context)
}

func (r Of) Go() (context.CancelFunc, error) {
	err := Spi.BeforeStart(&r)
	if err != nil {
		return nil, err
	}
	parent := r.ParentContext
	if parent == nil {
		parent = context.TODO()
	}
	ctx, cancel := context.WithCancel(parent)
	if r.OneOff != nil {
		r.goOneOff()
	} else {
		r.goLongRunning(ctx)
	}
	return cancel, nil
}

func (r *Of) goLongRunning(ctx context.Context) {
	go func() {
		for restartedTimes := 0; r.goLongRunningOnce(ctx); restartedTimes++ {
			shouldRestartAgain := Spi.BeforeRestart(r, restartedTimes)
			if !shouldRestartAgain {
				break
			}
		}
		Spi.AfterFinish(r)
	}()
}

func (r *Of) goLongRunningOnce(ctx context.Context) (notDone bool) {
	defer func() {
		recovered := recover()
		if recovered != nil {
			Spi.AfterPanic(r, recovered)
			notDone = true
		}
	}()
	r.LongRunning(ctx)
	return false
}

func (r *Of) goOneOff() {
	go func() {
		defer func() {
			recovered := recover()
			Spi.AfterPanic(r, recovered)
		}()
		r.OneOff()
		Spi.AfterFinish(r)
	}()
}

type Config struct {
	AfterPanic    func(routine *Of, recovered interface{})
	BeforeRestart func(routine *Of, restartedTimes int) bool
	BeforeStart   func(routine *Of) error
	AfterFinish   func(routine *Of)
}

var Spi = Config{
	AfterPanic: func(routine *Of, recovered interface{}) {
		// no op
	},
	BeforeRestart: func(routine *Of, restartedTimes int) bool {
		time.Sleep(100 * time.Microsecond)
		return true
	},
	BeforeStart: func(routine *Of) error {
		return nil
	},
	AfterFinish: func(routine *Of) {
		// no op
	},
}

func (cfg *Config) Append(newCfg Config) {
	if newCfg.AfterPanic != nil {
		oldAfterPanic := cfg.AfterPanic
		cfg.AfterPanic = func(routine *Of, recovered interface{}) {
			oldAfterPanic(routine, recovered)
			newCfg.AfterPanic(routine, recovered)
		}
	}
	if newCfg.BeforeRestart != nil {
		oldBeforeRestart := cfg.BeforeRestart
		cfg.BeforeRestart = func(routine *Of, restartedTimes int) bool {
			if !oldBeforeRestart(routine, restartedTimes) {
				return false
			}
			return newCfg.BeforeRestart(routine, restartedTimes)
		}
	}
	if newCfg.BeforeStart != nil {
		oldBeforeStart := cfg.BeforeStart
		cfg.BeforeStart = func(routine *Of) error {
			err := oldBeforeStart(routine)
			if err != nil {
				return err
			}
			return newCfg.BeforeStart(routine)
		}
	}
	if newCfg.AfterFinish != nil {
		oldAfterFinish := cfg.AfterFinish
		cfg.AfterFinish = func(routine *Of) {
			oldAfterFinish(routine)
			newCfg.AfterFinish(routine)
		}
	}
}
