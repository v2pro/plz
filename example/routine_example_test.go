package example

import (
	"fmt"
	"github.com/v2pro/plz"
	"github.com/v2pro/plz/routine"
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
	routine.Spi.Append(routine.Config{
		AfterPanic: func(recovered interface{}, kv []interface{}) {
			fmt.Println(recovered)
		},
	})
	defer func() {
		// restore back, after test
		routine.Spi.AfterPanic = func(recovered interface{}, kv []interface{}) {
		}
	}()
	plz.Go(func() {
		panic("hello")
	})
	time.Sleep(time.Second)
	// Output: hello
}
