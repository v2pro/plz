package must

import (
	"github.com/v2pro/plz/check"
	"runtime"
)

//go:noinline
func Check(result bool) {
	if result {
		return
	}
	t := check.CurrentT()
	t.Helper()
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		t.Fatal("check failed")
		return
	}
	t.Fatal(check.ExtractFailedLines(file, line))
}

//go:noinline
func Pass(result bool) {
	if result {
		return
	}
	t := check.CurrentT()
	t.Helper()
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		t.Fatal("check failed")
		return
	}
	t.Fatal(check.ExtractFailedLines(file, line))
}
