package countlog

import (
	"runtime"
	"fmt"
	"strings"
	"os"
	"time"
	"sync"
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
			level, event, expandedProperties = expand(level, event, properties)
		}
		defaultLogWriter.WriteLog(level, event, expandedProperties)
		return
	}
	for _, logWriter := range LogWriters {
		if !logWriter.ShouldLog(level, event, properties) {
			continue
		}
		if expandedProperties == nil {
			level, event, expandedProperties = expand(level, event, properties)
		}
		logWriter.WriteLog(level, event, expandedProperties)
	}
}
func expand(level int, event string, properties []interface{}) (int, string, []interface{}) {
	expandedProperties := make([]interface{}, 0, len(properties))
	hasError := false
	hasTimestamp := false
	for i := 0; i < len(properties); i = i + 2 {
		k := properties[i]
		v := properties[i+1]
		switch k {
		case "err":
			hasError = v != nil
		case "timestamp":
			hasTimestamp = true
		}
		expandedProperties = append(expandedProperties, k)
		switch typedProp := v.(type) {
		case func() interface{}:
			expandedProperties = append(expandedProperties, typedProp())
		case []byte:
			// []byte is likely being reused, need to make a copy here
			expandedProperties = append(expandedProperties, encodeAnyByteArray(typedProp))
		default:
			expandedProperties = append(expandedProperties, v)
		}
	}
	if !hasTimestamp {
		expandedProperties = append(expandedProperties, "timestamp")
		expandedProperties = append(expandedProperties, time.Now().UnixNano())
	}
	_, file, line, ok := runtime.Caller(3)
	if ok {
		expandedProperties = append(expandedProperties, "lineNumber")
		lineNumber := fmt.Sprintf("%s:%d", file, line)
		expandedProperties = append(expandedProperties, lineNumber)
		if event == "metric!" {
			if hasError {
				level = LevelError
			}
		} else if strings.HasPrefix(event, "event!") {
			event = event[len("event!"):]
		} else {
			os.Stderr.Write([]byte("countlog event must starts with event! prefix:" + lineNumber + "\n"))
			os.Stderr.Sync()
		}
	}
	return level, event, expandedProperties
}

// pull state callbacks

// like JMX MBean
type StateExporter interface {
	ExportState() map[string]interface{}
}

var stateExporters = map[string]StateExporter{}
var stateExportersMutex = &sync.Mutex{}

func RegisterStateExporter(name string, se StateExporter) {
	stateExportersMutex.Lock()
	defer stateExportersMutex.Unlock()
	stateExporters[name] = se
}

func RegisterStateExporterByFunc(name string, f func() map[string]interface{}) {
	stateExportersMutex.Lock()
	defer stateExportersMutex.Unlock()
	stateExporters[name] = &funcStateExporter{f}
}

func StateExporters() map[string]StateExporter {
	stateExportersMutex.Lock()
	defer stateExportersMutex.Unlock()
	return stateExporters
}

type funcStateExporter struct {
	f func() map[string]interface{}
}

func (se *funcStateExporter) ExportState() map[string]interface{} {
	return se.f()
}