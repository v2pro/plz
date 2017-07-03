package example

import (
	"fmt"
	"github.com/v2pro/plz"
	"github.com/v2pro/plz/lang/routine"
	"time"
)

func Example_go() {
	//routine.BeforeStart = append(routine.BeforeStart, func(kv []interface{}) error {
	//	return errors.New("can not start more goroutine")
	//})
	routine.AfterStart = append(routine.AfterStart, func(kv []interface{}) {
		fmt.Println("started")
	})
	routine.BeforeFinish = append(routine.BeforeFinish, func(kv []interface{}) {
		fmt.Println("finished")
	})
	plz.Go(func() {
		fmt.Println("hello from one off goroutine")
		panic("should not crash the whole program")
	})
	time.Sleep(time.Second)
	// Output:
	// started
	// hello from one off goroutine
	// finished
}

func Example_long_running_goroutine() {
	plz.GoLongRunning(func() {
		timer := time.NewTimer(time.Second).C
		for {
			select {
			case <-timer:
				fmt.Println("hello from running goroutine")
				return
			}
		}
	})
	time.Sleep(time.Second * 1)
	// Output: hello from running goroutine
}

func Example_log_goroutine_panic() {
	routine.AfterPanic = append(routine.AfterPanic, func(recovered interface{}, kv []interface{}) {
		fmt.Println("panic", recovered)
	})
	plz.Go(func() {
		panic("hello")
	})
	time.Sleep(time.Second)
	// Output: panic hello
}
