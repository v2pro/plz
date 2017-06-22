package example

import (
	"github.com/v2pro/plz"
	"fmt"
	"time"
	"context"
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
	cancel, _ := plz.Routine{LongRunning: func(ctx context.Context) bool {
		timer := time.NewTimer(time.Second).C
		for {
			select {
			case <-ctx.Done(): // to support cancel
				return false // return false, goroutine will not be restarted
			case <-timer:
				fmt.Println("hello from running goroutine")
			}
		}
		return true // if return true, goroutine will be restarted
	}}.Go()
	time.Sleep(time.Second * 2)
	cancel()
	// Output: hello from running goroutine
}
