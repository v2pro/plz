package countlog

import (
	"testing"
	"time"
	"errors"
)

func TestOsFileLogOutput(t *testing.T) {
	logWriter := NewAsyncLogWriter(LevelTrace, NewFileLogOutput("STDOUT"))
	logWriter.LogFormatter = &CompactFormat{StringLengthCap: 512}
	logWriter.Start()
	LogWriters = append(LogWriters, logWriter)
	Info("event!this is a test info")
}

func TestNewFileLogOutput(t *testing.T) {
	logWriter := NewAsyncLogWriter(LevelTrace, NewFileLogOutput("/tmp/test.log"))
	logWriter.LogFormatter = &CompactFormat{StringLengthCap: 512}
	logWriter.Start()
	LogWriters = append(LogWriters, logWriter)
	Info("event!this is a test info")
	time.Sleep(time.Second)
	logWriter.Close()
}

func Test_metric(t *testing.T) {
	Trace("metric!", "callee", "hello", "err", errors.New("my fault"))
}