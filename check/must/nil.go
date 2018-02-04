package must

import (
	"github.com/v2pro/plz/check"
	"github.com/v2pro/plz/check/testify/assert"
	"runtime"
)

func Nil(actual interface{}) {
	t := check.CurrentT()
	if assert.Nil(t, actual) {
		return
	}
	t.Helper()
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		t.Fatal("check failed")
		return
	}
	t.Fatal(check.ExtractFailedLines(file, line))
}

func AssertNil(actual interface{}) {
	t := check.CurrentT()
	if assert.Nil(t, actual) {
		return
	}
	t.Helper()
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		t.Fatal("check failed")
		return
	}
	t.Fatal(check.ExtractFailedLines(file, line))
}

func NotNil(actual interface{}) {
	t := check.CurrentT()
	if assert.NotNil(t, actual) {
		return
	}
	t.Helper()
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		t.Fatal("check failed")
		return
	}
	t.Fatal(check.ExtractFailedLines(file, line))
}

func AssertNotNil(actual interface{}) {
	t := check.CurrentT()
	if assert.NotNil(t, actual) {
		return
	}
	t.Helper()
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		t.Fatal("check failed")
		return
	}
	t.Fatal(check.ExtractFailedLines(file, line))
}