package countlog

import (
	"testing"
	"github.com/v2pro/plz/countlog/compact"
	"io/ioutil"
)

func Test_trace(t *testing.T) {
	MinLevel = LevelTrace
	Trace("event!hello", "a", "b")
}

func Benchmark_trace(b *testing.B) {
	DefaultEventSink = &WriteEventSink{
		Format: &compact.Format{},
		Writer: ioutil.Discard,
	}
	MinLevel = LevelTrace
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Trace("event!hello", "a", "b")
	}
}
