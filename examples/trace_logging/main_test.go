package main

import (
	"testing"
	"github.com/v2pro/plz/countlog"
)

func Benchmark_trace_call(b *testing.B) {
	for i := 0; i < b.N; i++ {
		countlog.TraceCall("callee!doSomething", nil)
	}
}