package countlog

import (
	"runtime"
	"fmt"
	"strings"
	"os"
	"time"
)

// push event out

const LevelTrace = 10
const LevelDebug = 20
const LevelInfo = 30
const LevelWarn = 40
const LevelError = 50
const LevelFatal = 60

func getLevelName(level int) string {
	switch level {
	case LevelTrace: return "TRACE"
	case LevelDebug: return "DEBUG"
	case LevelInfo: return "INFO"
	case LevelWarn: return "WARN"
	case LevelError: return "ERROR"
	case LevelFatal: return "FATAL"
	default: return "UNKNOWN"
	}
}

func Trace(event string, properties ...interface{}) {
	log(LevelTrace, event, properties)
}
func Debug(event string, properties ...interface{}) {
	log(LevelDebug, event, properties)
}
func Info(event string, properties ...interface{}) {
	log(LevelInfo, event, properties)
}
func Warn(event string, properties ...interface{}) {
	log(LevelWarn, event, properties)
}
func Error(event string, properties ...interface{}) {
	log(LevelError, event, properties)
}
func Fatal(event string, properties ...interface{}) {
	log(LevelFatal, event, properties)
}
func Log(level int, event string, properties ...interface{}) {
	log(level, event, properties)
}
func log(level int, event string, properties []interface{}) {
	var expandedProperties []interface{}
	if len(LogWriters) == 0 {
		if expandedProperties == nil {
			event, expandedProperties = expand(event, properties)
		}
		defaultLogWriter.WriteLog(level, event, expandedProperties)
		return
	}
	for _, logWriter := range LogWriters {
		if !logWriter.ShouldLog(level, event, properties) {
			continue
		}
		if expandedProperties == nil {
			event, expandedProperties = expand(event, properties)
		}
		logWriter.WriteLog(level, event, expandedProperties)
	}
}
func expand(event string, properties []interface{}) (string, []interface{}) {
	expandedProperties := []interface{}{
		"timestamp", time.Now().UnixNano(),
	}
	_, file, line, ok := runtime.Caller(3)
	if ok {
		expandedProperties = append(expandedProperties, "lineNumber")
		lineNumber := fmt.Sprintf("%s:%d", file, line)
		expandedProperties = append(expandedProperties, lineNumber)
		if strings.HasPrefix(event, "event!") {
			event = event[len("event!"):]
		} else {
			os.Stderr.Write([]byte("countlog event must starts with event! prefix:" + lineNumber + "\n"))
			os.Stderr.Sync()
		}
	}
	for _, prop := range properties {
		switch typedProp := prop.(type) {
		case func() interface{}:
			expandedProperties = append(expandedProperties, typedProp())
		case []byte:
			// []byte is likely being reused, need to make a copy here
			expandedProperties = append(expandedProperties, encodeAnyByteArray(typedProp))
		default:
			expandedProperties = append(expandedProperties, prop)
		}
	}
	return event, expandedProperties
}

// pull state callbacks

// like JMX MBean
type StateExporter interface {
	ExportState() map[string]interface{}
}

var StateExporters = map[string]StateExporter{}