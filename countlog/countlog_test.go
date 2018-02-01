package countlog

import (
	"testing"
	"github.com/v2pro/plz/countlog/compact"
	"io/ioutil"
	"time"
	"os"
)

func Test_trace(t *testing.T) {
	MinLevel = LevelTrace
	executor := &defaultExecutor{}
	DevelopmentEventSink = NewEventWriter(EventWriterConfig{
		Format:   &compact.Format{},
		Writer:   os.Stdout,
		Executor: executor,
	})
	Trace("event!hello", "a", "b", "int", 100)
	time.Sleep(time.Second)
}

func Test_trace_call(t *testing.T) {
	MinLevel = LevelTrace
	TraceCall("callee!func", nil)
}

func Benchmark_trace(b *testing.B) {
	DevelopmentEventSink = NewEventWriter(EventWriterConfig{
		Format: &compact.Format{},
		Writer: ioutil.Discard,
	})
	MinLevel = LevelTrace
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Trace("event!hello", "a", "b", "int", 100)
	}
}
