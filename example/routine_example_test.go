package example

import (
	"fmt"
	"github.com/v2pro/plz"
	"github.com/v2pro/plz/lang/routine"
	"time"
)

func Example_go() {
	plz.Go(func() {
		fmt.Println("hello from one off goroutine")
		panic("should not crash the whole program")
	})
	time.Sleep(time.Second)
	// Output: hello from one off goroutine
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
