package plz

import (
	"testing"
	"github.com/json-iterator/go/require"
	"sync"
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
