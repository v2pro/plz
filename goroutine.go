package plz

import (
	"context"
	"time"
)

func Go(oneOff func()) error {
	_, err := Routine{OneOff: oneOff}.Go()
	return err
}

type Routine struct {
	ParentContext context.Context
	OneOff        func()
	LongRunning   func(ctx context.Context) bool
}

func (r Routine) Go() (context.CancelFunc, error) {
	err := RoutineSpi.BeforeStart(&r)
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

func (r *Routine) goLongRunning(ctx context.Context) {
	go func() {
		for restartedTimes := 0; r.goLongRunningOnce(ctx); restartedTimes++ {
			shouldRestartAgain := RoutineSpi.BeforeRestart(r, restartedTimes)
			if !shouldRestartAgain {
				break
			}
		}
		RoutineSpi.AfterFinish(r)
	}()
}

func (r *Routine) goLongRunningOnce(ctx context.Context) (notDone bool) {
	defer func() {
		recovered := recover()
		if recovered != nil {
			RoutineSpi.AfterPanic(r, recovered)
			notDone = true
		}
	}()
	return r.LongRunning(ctx)
}

func (r *Routine) goOneOff() {
	go func() {
		defer func() {
			recovered := recover()
			RoutineSpi.AfterPanic(r, recovered)
		}()
		r.OneOff()
		RoutineSpi.AfterFinish(r)
	}()
}

type RoutineSpiConfig struct {
	AfterPanic    func(routine *Routine, recovered interface{})
	BeforeRestart func(routine *Routine, restartedTimes int) bool
	BeforeStart   func(routine *Routine) error
	AfterFinish   func(routine *Routine)
}

var RoutineSpi = RoutineSpiConfig{
	AfterPanic: func(routine *Routine, recovered interface{}) {
		// no op
	},
	BeforeRestart: func(routine *Routine, restartedTimes int) bool {
		time.Sleep(100 * time.Microsecond)
		return true
	},
	BeforeStart: func(routine *Routine) error {
		return nil
	},
	AfterFinish: func(routine *Routine) {
		// no op
	},
}

func (cfg *RoutineSpiConfig) Append(newCfg RoutineSpiConfig) {
	if newCfg.AfterPanic != nil {
		oldAfterPanic := cfg.AfterPanic
		cfg.AfterPanic = func(routine *Routine, recovered interface{}) {
			oldAfterPanic(routine, recovered)
			newCfg.AfterPanic(routine, recovered)
		}
	}
	if newCfg.BeforeRestart != nil {
		oldBeforeRestart := cfg.BeforeRestart
		cfg.BeforeRestart = func(routine *Routine, restartedTimes int) bool {
			if !oldBeforeRestart(routine, restartedTimes) {
				return false
			}
			return newCfg.BeforeRestart(routine, restartedTimes)
		}
	}
	if newCfg.BeforeStart != nil {
		oldBeforeStart := cfg.BeforeStart
		cfg.BeforeStart = func(routine *Routine) error {
			err := oldBeforeStart(routine)
			if err != nil {
				return err
			}
			return newCfg.BeforeStart(routine)
		}
	}
	if newCfg.AfterFinish != nil {
		oldAfterFinish := cfg.AfterFinish
		cfg.AfterFinish = func(routine *Routine) {
			oldAfterFinish(routine)
			newCfg.AfterFinish(routine)
		}
	}
}
