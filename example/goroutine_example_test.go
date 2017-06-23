package example

import (
	"github.com/v2pro/plz"
	"fmt"
	"time"
	"context"
	"github.com/v2pro/plz/routine"
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
	cancel, _ := plz.GoLongRunning(func(ctx context.Context) {
		timer := time.NewTimer(time.Second).C
		for {
			select {
			case <-ctx.Done(): // to support cancel
				return
			case <-timer:
				fmt.Println("hello from running goroutine")
				return
			}
		}
	})
	time.Sleep(time.Second * 2)
	cancel()
	// Output: hello from running goroutine
}

func Example_log_goroutine_panic() {
	routine.Spi.Append(routine.Config{
		AfterPanic: func(routine *routine.Of, recovered interface{}) {
			fmt.Println(recovered)
		},
	})
	defer func() {
		// restore back, after test
		routine.Spi.AfterPanic = func(routine *routine.Of, recovered interface{}) {
		}
	}()
	routine.Of{OneOff: func() { // same as plz.Go()
		panic("hello")
	}}.Go()
	time.Sleep(time.Second)
	// Output: hello
}
