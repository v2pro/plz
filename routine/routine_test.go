package routine

import (
	"testing"
	"github.com/json-iterator/go/require"
	"sync"
	"context"
	"time"
	"errors"
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
	RoutineSpi.AfterPanic = func(routine *Routine, recovered interface{}) {
		lock.Unlock()
		called = true
	}
	defer func() {
		RoutineSpi.AfterPanic = func(routine *Routine, recovered interface{}) {
		}
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
	Routine{LongRunning: func(ctx context.Context) {
		counter++
		if counter > 3 {
			lock.Unlock()
		}
		panic("hello")
	}}.Go()
	lock.Lock()
	should.Equal(4, counter)
}

func Test_long_running_goroutine_cancel(t *testing.T) {
	should := require.New(t)
	called := false
	lock := &sync.Mutex{}
	lock.Lock()
	cancel, _ := Routine{LongRunning: func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				lock.Unlock()
				called = true
			default:
				time.Sleep(time.Second)
			}
		}
	}}.Go()
	cancel()
	lock.Lock()
	should.True(called)
}

func Test_routine_spi_before_start(t *testing.T) {
	should := require.New(t)
	RoutineSpi.BeforeStart = func(routine *Routine) error {
		return errors.New("exceed limit")
	}
	defer func() {
		RoutineSpi.BeforeStart = func(routine *Routine) error {
			return nil
		}
	}()
	should.NotNil(Go(func() {}))
}

func Test_routine_spi_after_finish(t *testing.T) {
	should := require.New(t)
	called := false
	lock := &sync.Mutex{}
	RoutineSpi.AfterFinish = func(routine *Routine) {
		called = true
		lock.Unlock()
	}
	defer func() {
		RoutineSpi.AfterFinish = func(routine *Routine) {
		}
	}()
	lock.Lock()
	Go(func() {})
	lock.Lock()
	should.True(called)
}

func Test_routine_spi_composition(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(2)
	RoutineSpi.Append(RoutineSpiConfig{
		AfterFinish: func(routine *Routine) {
			wg.Done()
		},
	})
	RoutineSpi.Append(RoutineSpiConfig{
		AfterFinish: func(routine *Routine) {
			wg.Done()
		},
	})
	defer func() {
		RoutineSpi.AfterFinish = func(routine *Routine) {
		}
	}()
	Go(func() {})
	wg.Wait()
}
