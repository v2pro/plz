package should

import (
	"github.com/v2pro/plz/check"
	"runtime"
)

//go:noinline
func Check(result bool) {
	if !result {
		t := check.CurrentT()
		t.Helper()
		_, file, line, ok := runtime.Caller(1)
		if !ok {
			t.Error("check failed")
			return
		}
		t.Error(check.ExtractFailedLines(file, line))
	}
}

//go:noinline
func Pass(result bool) {
	if !result {
		t := check.CurrentT()
		t.Helper()
		_, file, line, ok := runtime.Caller(1)
		if !ok {
			t.Error("check failed")
			return
		}
		t.Error(check.ExtractFailedLines(file, line))
	}
}
