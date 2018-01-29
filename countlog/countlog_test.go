package countlog

import (
	"errors"
	"testing"
	"time"
)

func TestOsFileLogOutput(t *testing.T) {
	logWriter := NewAsyncLogWriter(LevelTrace, NewFileLogOutput("STDOUT"))
	logWriter.LogFormat = &CompactFormat{StringLengthCap: 512}
	logWriter.Start()
	LogWriters = append(LogWriters, logWriter)
	Info("event!this is a test info")
}

func TestNewFileLogOutput(t *testing.T) {
	logWriter := NewAsyncLogWriter(LevelTrace, NewFileLogOutput("/tmp/test.log"))
	logWriter.LogFormat = &CompactFormat{StringLengthCap: 512}
	logWriter.Start()
	LogWriters = append(LogWriters, logWriter)
	Info("event!this is a test info")
	time.Sleep(time.Second)
	logWriter.Close()
}

func Test_metric(t *testing.T) {
	err := errors.New("my fault")
	// when performance is critical, use ShouldLog to reduce log overhead
	if err != nil || ShouldLog(LevelTrace) {
		TraceCall(err, "callee", "hello")
	}
	// without ShouldLog the overhead is still minimum
	TraceCall(err, "callee", "hello")
	// err == nil will not show, because Trace < Debug
	TraceCall(nil, "callee", "world")
}
