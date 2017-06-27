package routine

import (
	"errors"
	"github.com/json-iterator/go/require"
	"sync"
	"testing"
)

func Test_one_off_goroutine_long_version(t *testing.T) {
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
	Spi.AfterPanic = func(recovered interface{}, kv []interface{}) {
		lock.Unlock()
		called = true
	}
	defer func() {
		Spi.AfterPanic = func(recovered interface{}, kv []interface{}) {
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
	GoLongRunning(func() {
		counter++
		if counter > 3 {
			lock.Unlock()
		}
		panic("hello")
	})
	lock.Lock()
	should.Equal(4, counter)
}

func Test_routine_spi_before_start(t *testing.T) {
	should := require.New(t)
	Spi.BeforeStart = func(kv []interface{}) error {
		return errors.New("exceed limit")
	}
	defer func() {
		Spi.BeforeStart = func(kv []interface{}) error {
			return nil
		}
	}()
	should.NotNil(Go(func() {}))
}

func Test_routine_spi_after_finish(t *testing.T) {
	should := require.New(t)
	called := false
	lock := &sync.Mutex{}
	Spi.AfterFinish = func(kv []interface{}) {
		called = true
		lock.Unlock()
	}
	defer func() {
		Spi.AfterFinish = func(kv []interface{}) {
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
		AfterFinish: func(kv []interface{}) {
			wg.Done()
		},
	})
	Spi.Append(Config{
		AfterFinish: func(kv []interface{}) {
			wg.Done()
		},
	})
	defer func() {
		Spi.AfterFinish = func(kv []interface{}) {
		}
	}()
	Go(func() {})
	wg.Wait()
}
