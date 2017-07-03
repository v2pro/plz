package routine

import (
	"fmt"
	"github.com/v2pro/plz/logging"
	"runtime"
	"runtime/debug"
	"time"
)

var panicLogger = logging.LoggerOf("metric", "counter", "panic", "routine")

func Go(oneOff func(), kv ...interface{}) error {
	var err error
	for _, handle := range BeforeStart {
		err = handle(kv)
		if err != nil {
			return err
		}
	}
	_, callerFile, callerLine, _ := runtime.Caller(2)
	go func() {
		defer func() {
			recovered := recover()
			if recovered != nil {
				for _, handle := range AfterPanic {
					handle(recovered, append(kv, "caller", fmt.Sprintf("%s:%d", callerFile, callerLine)))
				}
			}
			for _, handle := range BeforeFinish {
				handle(kv)
			}
		}()
		oneOff()
	}()
	return nil
}

func GoLongRunning(longRunning func(), kv ...interface{}) error {
	var err error
	for _, handle := range BeforeStart {
		err = handle(kv)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	go func() {
		defer func() {
			for _, handle := range BeforeFinish {
				handle(kv)
			}
		}()
		for restartedTimes := 0; goLongRunningOnce(longRunning, kv); restartedTimes++ {
			for _, handle := range BeforeRestart {
				if !handle(restartedTimes, kv) {
					return
				}
			}
		}
	}()
	return nil
}

func goLongRunningOnce(longRunning func(), kv []interface{}) (notDone bool) {
	defer func() {
		recovered := recover()
		if recovered != nil {
			for _, handle := range AfterPanic {
				handle(recovered, kv)
			}
			notDone = true
		}
	}()
	longRunning()
	return false
}

type Config struct {
	AfterPanic    func(recovered interface{}, kv []interface{})
	BeforeRestart func(restartedTimes int, kv []interface{}) bool
	BeforeStart   func(kv []interface{}) error
	AfterFinish   func(kv []interface{})
}

var AfterPanic = []func(recovered interface{}, kv []interface{}){
	func(recovered interface{}, kv []interface{}) {
		panicLogger.Error("goroutine panic", append(kv, "recovered", recovered, "stack", string(debug.Stack()))...)
	},
}

var BeforeRestart = []func(restartedTimes int, kv []interface{}) bool{
	func(restartedTimes int, kv []interface{}) bool {
		time.Sleep(100 * time.Microsecond)
		return true
	},
}

var BeforeStart = []func(kv []interface{}) error{
	func(kv []interface{}) error {
		return nil // allow go without limit
	},
}

var BeforeFinish = []func(kv []interface{}){
	func(kv []interface{}) {
		// no op
	},
}
