package countlog

import (
	"testing"
	"time"
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"github.com/v2pro/plz/countlog/output"
	"github.com/v2pro/plz/countlog/output/compact"
	"os"
	"github.com/v2pro/plz/countlog/output/lumberjack"
)

func Test_trace(t *testing.T) {
	Trace("event!hello", "a", "b", "int", 100)
	time.Sleep(time.Second)
}

func Test_trace_call(t *testing.T) {
	should := require.New(t)
	err := DebugCall("call func with %(k1)s", errors.New("failure"),
		"k1", "v1")
	should.Equal("call func with v1: failure", err.Error())
}

func Test_call_with_same_event_but_different_properties(t *testing.T) {
	ctx := Ctx(context.Background())
	for i := 0; i < 3; i++ {
		ctx.Trace("same event name", "key", 100)
		Trace("same event name", "key", "value")
	}
}

func Test_log_file(t *testing.T) {
	should := require.New(t)
	logFile, err := os.OpenFile("/tmp/test.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	should.NoError(err)
	defer logFile.Close()
	EventWriter = output.NewEventWriter(output.EventWriterConfig{
		Format: &compact.Format{},
		Writer: output.NewAsyncWriter(output.AsyncWriterConfig{
			QueueLength: 1024,
			IsQueueBlocking: false,
			Writer: logFile,
		}),
	})
	for i := 0; i < 1000; i++ {
		Info("something happened", "input", "abc", "output", "def")
	}
	time.Sleep(time.Second)
}

func Test_rolling_log_file(t *testing.T) {
	logFile := &lumberjack.Logger{
		BackupTimeFormat: "2006-01-02T15-04-05.000",
		Filename:   "/tmp/test.log",
		MaxSize:    1, // megabytes
		MaxBackups: 3,
	}
	defer logFile.Close()
	EventWriter = output.NewEventWriter(output.EventWriterConfig{
		Format: &compact.Format{},
		Writer: logFile,
	})
	for i := 0; i < 10000; i++ {
		Info("something happened", "input", "abc", "output", "def")
	}
}

func Benchmark_trace(b *testing.B) {
	SetMinLevel(LevelDebug)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		TraceCall("event!hello", nil, "a", "b")
	}
}
