package countlog

import (
	"runtime"
	"fmt"
	"strings"
	"os"
	"time"
)

const LEVEL_TRACE = 10
const LEVEL_DEBUG = 20
const LEVEL_INFO = 30
const LEVEL_WARN = 40
const LEVEL_ERROR = 50
const LEVEL_FATAL = 60

func Trace(event string, properties ...interface{}) {
	log(LEVEL_TRACE, event, properties)
}
func Debug(event string, properties ...interface{}) {
	log(LEVEL_DEBUG, event, properties)
}
func Info(event string, properties ...interface{}) {
	log(LEVEL_INFO, event, properties)
}
func Warn(event string, properties ...interface{}) {
	log(LEVEL_WARN, event, properties)
}
func Error(event string, properties ...interface{}) {
	log(LEVEL_ERROR, event, properties)
}
func Fatal(event string, properties ...interface{}) {
	log(LEVEL_FATAL, event, properties)
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
		propProvider, _ := prop.(func() interface{})
		if propProvider == nil {
			expandedProperties = append(expandedProperties, prop)
		} else {
			expandedProperties = append(expandedProperties, propProvider())
		}
	}
	return event, expandedProperties
}
