package countlog

import (
	"testing"
	"time"
	"os"
	"github.com/v2pro/plz/countlog/output"
	"github.com/v2pro/plz/countlog/output/compact"
	"io/ioutil"
)

func Test_trace(t *testing.T) {
	DevelopmentEventSink = output.NewEventWriter(output.EventWriterConfig{
		Format:   &compact.Format{},
		Writer:   os.Stdout,
		Executor: output.DefaultExecutor,
	})
	Trace("event!hello", "a", "b", "int", 100)
	time.Sleep(time.Second)
}

func Test_trace_call(t *testing.T) {
	DebugCall("callee!func", nil, "k1", "v1")
}

func Test_call_with_same_event_but_different_properties(t *testing.T) {
	Trace("same event name", "key", "value")
	Trace("same event name", "key", 100)
}

func Benchmark_trace(b *testing.B) {
	DevelopmentEventSink = output.NewEventWriter(output.EventWriterConfig{
		Format:   &compact.Format{},
		Writer:   ioutil.Discard,
		//Executor: output.DefaultExecutor,
	})
	SetMinLevel(LevelDebug)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		TraceCall("event!hello",nil, "a", "b")
	}
}
