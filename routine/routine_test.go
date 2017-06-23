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
	Of{OneOff: func() {
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
	Spi.AfterPanic = func(routine *Of, recovered interface{}) {
		lock.Unlock()
		called = true
	}
	defer func() {
		Spi.AfterPanic = func(routine *Of, recovered interface{}) {
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
	Of{LongRunning: func(ctx context.Context) {
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
	cancel, _ := Of{LongRunning: func(ctx context.Context) {
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
	Spi.BeforeStart = func(routine *Of) error {
		return errors.New("exceed limit")
	}
	defer func() {
		Spi.BeforeStart = func(routine *Of) error {
			return nil
		}
	}()
	should.NotNil(Go(func() {}))
}

func Test_routine_spi_after_finish(t *testing.T) {
	should := require.New(t)
	called := false
	lock := &sync.Mutex{}
	Spi.AfterFinish = func(routine *Of) {
		called = true
		lock.Unlock()
	}
	defer func() {
		Spi.AfterFinish = func(routine *Of) {
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
	Spi.Append(Config{
		AfterFinish: func(routine *Of) {
			wg.Done()
		},
	})
	Spi.Append(Config{
		AfterFinish: func(routine *Of) {
			wg.Done()
		},
	})
	defer func() {
		Spi.AfterFinish = func(routine *Of) {
		}
	}()
	Go(func() {})
	wg.Wait()
}
