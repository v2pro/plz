package plz

import (
	"testing"
	"github.com/json-iterator/go/require"
	"sync"
	"context"
	"time"
)

func Test_one_off_goroutine_long_version(t *testing.T) {
	should := require.New(t)
	called := false
	lock := &sync.Mutex{}
	lock.Lock()
	Routine{OneOff: func() {
		lock.Unlock()
		called = true
	}}.Go()
	lock.Lock()
	should.True(called)
}

func Test_one_off_goroutine_short_version(t *testing.T) {
	should := require.New(t)
	called := false
	lock := &sync.Mutex{}
	lock.Lock()
	Go(func() {
		lock.Unlock()
		called = true
	})
	lock.Lock()
	should.True(called)
}

func Test_one_off_goroutine_panic(t *testing.T) {
	should := require.New(t)
	called := false
	lock := &sync.Mutex{}
	lock.Lock()
	RegisterPanicHandler(func(routine *Routine, recovered interface{}) {
		lock.Unlock()
		called = true
	})
	defer func() {
		panicHandlers = []HandlePanic{}
	}()
	Go(func() {
		panic("hello")
	})
	lock.Lock()
	should.True(called)
}

func Test_long_running_goroutine_should_be_restarted(t *testing.T) {
	should := require.New(t)
	counter := 0
	lock := &sync.Mutex{}
	lock.Lock()
	Routine{LongRunning: func(ctx context.Context) bool {
		counter++
		if counter > 3 {
			lock.Unlock()
			return false
		}
		panic("hello")
		return true
	}}.Go()
	lock.Lock()
	should.Equal(4, counter)
}

func Test_long_running_goroutine_cancel(t *testing.T) {
	should := require.New(t)
	called := false
	lock := &sync.Mutex{}
	lock.Lock()
	cancel := Routine{LongRunning: func(ctx context.Context) bool {
		for {
			select {
			case <-ctx.Done():
				lock.Unlock()
				called = true
				return false
			default:
				time.Sleep(time.Second)
			}
		}
		return true
	}}.Go()
	cancel()
	lock.Lock()
	should.True(called)
}
