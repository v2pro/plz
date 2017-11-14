package countlog

import (
	"testing"
)

func TestNewFileLogOutput(t *testing.T) {
	logWriter := NewAsyncLogWriter(LevelTrace, NewFileLogOutput("STDOUT"))
	logWriter.LogFormatter = &CompactFormat{StringLengthCap: 512}
	logWriter.Start()
	LogWriters = append(LogWriters, logWriter)
	Info("event!this is a test info")
}