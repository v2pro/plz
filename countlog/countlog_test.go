package countlog

import (
	"testing"
	"time"
	"os"
	"github.com/v2pro/plz/countlog/output"
	"github.com/v2pro/plz/countlog/output/compact"
)

func Test_trace(t *testing.T) {
	MinLevel = LevelTrace
	DevelopmentEventSink = output.NewEventWriter(output.EventWriterConfig{
		Format:   &compact.Format{},
		Writer:   os.Stdout,
		Executor: output.DefaultExecutor,
	})
	Trace("event!hello", "a", "b", "int", 100)
	time.Sleep(time.Second)
}

func Test_trace_call(t *testing.T) {
	MinLevel = LevelTrace
	TraceCall("callee!func", nil)
}

func Benchmark_trace(b *testing.B) {
	DevelopmentEventSink = output.NewEventWriter(output.EventWriterConfig{
		Format:   &compact.Format{},
		Writer:   os.Stdout,
		Executor: output.DefaultExecutor,
	})
	MinLevel = LevelTrace
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Trace("event!hello", "a", "b", "int", 100)
	}
}
