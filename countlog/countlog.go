package countlog

import (
	"sync"
	"unsafe"
	"runtime"
	"github.com/v2pro/plz/countlog/core"
	"time"
)

// push event out

const LevelTrace = 10
const LevelDebug = 20
const LevelInfo = 30
const LevelWarn = 40
const LevelError = 50
const LevelFatal = 60

// MinLevel exists to minimize the overhead of Trace/Debug logging
var MinLevel = LevelDebug

func getLevelName(level int) string {
	switch level {
	case LevelTrace:
		return "TRACE"
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

func ShouldLog(level int) bool {
	return level >= MinLevel
}

func Trace(event string, properties ...interface{}) {
	if LevelTrace < MinLevel {
		return
	}
	ptr := unsafe.Pointer(&properties)
	log(LevelTrace, event, nil, nil, castEmptyInterfaces(uintptr(ptr)))
}

func castEmptyInterfaces(ptr uintptr) []interface{} {
	return *(*[]interface{})(unsafe.Pointer(ptr))
}

func TraceCall(callee string, err error, properties ...interface{}) {
	if err != nil {
		log(LevelError, callee, nil, err, properties)
		return
	}
	if LevelTrace < MinLevel {
		return
	}
	log(LevelTrace, callee, nil, err, properties)
}

func Debug(event string, properties ...interface{}) {
	if LevelDebug < MinLevel {
		return
	}
	log(LevelDebug, event, nil, nil, properties)
}

func DebugCall(callee string, err error, properties ...interface{}) {
	if err != nil {
		log(LevelError, callee, nil, err, properties)
		return
	}
	if LevelDebug < MinLevel {
		return
	}
	log(LevelDebug, callee, nil, err, properties)
}

func Info(event string, properties ...interface{}) {
	if LevelInfo < MinLevel {
		return
	}
	log(LevelInfo, event, nil, nil, properties)
}

func InfoCall(callee string, err error, properties ...interface{}) {
	if err != nil {
		log(LevelError, callee, nil, err, properties)
		return
	}
	if LevelInfo < MinLevel {
		return
	}
	log(LevelInfo, callee, nil, err, properties)
}

func Warn(event string, properties ...interface{}) {
	log(LevelWarn, event, nil, nil, properties)
}

func Error(event string, properties ...interface{}) {
	log(LevelError, event, nil, nil, properties)
}

func Fatal(event string, properties ...interface{}) {
	log(LevelFatal, event, nil, nil, properties)
}

func Log(level int, event string, properties ...interface{}) {
	log(level, event, nil, nil, properties)
}

func LogPanic(recovered interface{}, properties ...interface{}) interface{} {
	if recovered != nil {
		buf := make([]byte, 1<<16)
		runtime.Stack(buf, false)
		if len(properties) > 0 {
			properties = append(properties, "err", recovered, "stacktrace", string(buf))
			Fatal("event!panic", properties...)
		} else {
			Fatal("event!panic", "err", recovered, "stacktrace", string(buf))
		}
	}
	return recovered
}

var handlerCache = &sync.Map{}

func log(level int, eventOrCallee string, ctx *Context, err error, properties []interface{}) {
	handler := getHandler(level, eventOrCallee, ctx, properties)
	event := &core.Event{
		Context:    ctx,
		Error:      err,
		Timestamp:  time.Now(),
		Properties: properties,
	}
	ptr := unsafe.Pointer(event)
	handler.Handle(castEvent(uintptr(ptr)))
}

func castEvent(ptr uintptr) *core.Event {
	return (*core.Event)(unsafe.Pointer(ptr))
}
func castString(ptr uintptr) string {
	return *(*string)(unsafe.Pointer(ptr))
}

func getHandler(level int, eventOrCallee string, ctx *Context, properties []interface{}) core.EventHandler {
	handler, found := handlerCache.Load(eventOrCallee)
	if found {
		return handler.(core.EventHandler)
	}
	return newHandler(level, eventOrCallee, ctx, properties)
}

func newHandler(level int, eventOrCalleeObj string, ctx *Context, properties []interface{}) core.EventHandler {
	ptr := unsafe.Pointer(&eventOrCalleeObj)
	eventOrCallee := castString(uintptr(ptr))
	skipFramesCount := 3
	if ctx != nil {
		skipFramesCount = 5
	}
	_, callerFile, callerLine, _ := runtime.Caller(skipFramesCount)
	var handlers core.EventHandlers
	for _, sink := range EventSinks {
		handler := sink.HandlerOf(level, eventOrCallee, callerFile, callerLine, properties)
		if handler == nil {
			continue
		}
		handlers = append(handlers, handler)
	}
	switch len(handlers) {
	case 0:
		handler := DevelopmentEventSink.HandlerOf(level, eventOrCallee, callerFile, callerLine, properties)
		handlerCache.Store(eventOrCallee, handler)
		return handler
	case 1:
		handlerCache.Store(eventOrCallee, handlers[0])
		return handlers[0]
	default:
		handlerCache.Store(eventOrCallee, handlers)
		return handlers
	}
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
