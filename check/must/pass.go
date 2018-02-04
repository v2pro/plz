package must

import (
	"github.com/v2pro/plz/check"
	"runtime"
	"github.com/v2pro/plz/check/go-spew/spew"
)

//go:noinline
func Assert(result bool, kv ...interface{}) {
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
func Pass(result bool, kv ...interface{}) {
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
	for i := 0; i < len(kv); i+=2 {
		key := kv[i].(string)
		t.Errorf("%s: %s", key, spew.Sdump(kv[i+1]))
	}
	t.Fatal(check.ExtractFailedLines(file, line))
}
