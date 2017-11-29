package concurrent

import (
	"sync"
	"github.com/v2pro/plz/countlog"
	"context"
	"runtime"
	"time"
)

const StopSignal = "STOP!"

type UnboundedExecutor struct {
	ctx                   context.Context
	cancel                context.CancelFunc
	activeGoroutinesMutex *sync.Mutex
	activeGoroutines      map[startFrom]int
}

type startFrom struct {
	startFromFile string
	startFromLine int
}

func NewUnboundedExecutor() *UnboundedExecutor {
	ctx, cancel := context.WithCancel(context.TODO())
	return &UnboundedExecutor{
		ctx:                   ctx,
		cancel:                cancel,
		activeGoroutinesMutex: &sync.Mutex{},
		activeGoroutines:      map[startFrom]int{},
	}
}

func (executor *UnboundedExecutor) Go(handler func(ctx context.Context)) {
	_, file, line, _ := runtime.Caller(1)
	countlog.Info("event!unbounded_executor.start goroutine",
		"startFromFile", file, "startFromLine", line)
	executor.activeGoroutinesMutex.Lock()
	defer executor.activeGoroutinesMutex.Unlock()
	startFrom := startFrom{
		startFromFile: file,
		startFromLine: line,
	}
	executor.activeGoroutines[startFrom] += 1
	go func() {
		defer func() {
			recovered := recover()
			if recovered != nil && recovered != StopSignal {
				countlog.Fatal("event!unbounded_executor.panic",
					"err", recovered,
					"stacktrace", countlog.ProvideStacktrace)
			}
			executor.activeGoroutinesMutex.Lock()
			defer executor.activeGoroutinesMutex.Unlock()
			executor.activeGoroutines[startFrom] -= 1
		}()
		handler(executor.ctx)
	}()
}

func (executor *UnboundedExecutor) StopAndWait() {
	executor.cancel()
	for {
		if executor.checkGoroutines() {
			return
		}
		time.Sleep(time.Second * 3)
	}
}

func (executor *UnboundedExecutor) checkGoroutines() bool {
	executor.activeGoroutinesMutex.Lock()
	defer executor.activeGoroutinesMutex.Unlock()
	allDone := true
	for startFrom, count := range executor.activeGoroutines {
		if count > 0 {
			countlog.Info("event!unbounded_executor.still waiting goroutines to quit",
				"startFromFile", startFrom.startFromFile,
				"startFromLine", startFrom.startFromLine,
				"count", count)
			allDone = false
		}
	}
	return allDone
}
