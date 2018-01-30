package countlog

import (
	"testing"
)

func Test_trace(t *testing.T) {
	MinLevel = LevelTrace
	Trace("event!hello", "a", "b")
}
//
//func TestOsFileLogOutput(t *testing.T) {
//	logWriter := NewAsyncLogWriter(LevelTrace, NewFileLogOutput("STDOUT"))
//	logWriter.LogFormat = &CompactFormat{StringLengthCap: 512}
//	logWriter.Start()
//	LogWriters = append(LogWriters, logWriter)
//	Info("event!this is a test info")
//}
//
//func TestNewFileLogOutput(t *testing.T) {
//	logWriter := NewAsyncLogWriter(LevelTrace, NewFileLogOutput("/tmp/test.log"))
//	logWriter.LogFormat = &CompactFormat{StringLengthCap: 512}
//	logWriter.Start()
//	LogWriters = append(LogWriters, logWriter)
//	Info("event!this is a test info")
//	time.Sleep(time.Second)
//	logWriter.Close()
//}
