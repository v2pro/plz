package must

import (
	"github.com/stretchr/testify/assert"
	"github.com/v2pro/plz/check"
	"runtime"
)

//go:noinline
func Equal(expected interface{}, actual interface{}) {
	t := check.CurrentT()
	if assert.Equal(t, expected, actual) {
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

//go:noinline
func CheckEqual(expected interface{}, actual interface{}) {
	t := check.CurrentT()
	if assert.Equal(t, expected, actual) {
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
