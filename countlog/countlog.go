package countlog

import (
	"sync"
	"unsafe"
	"runtime"
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

var handlerCache = &sync.Map{}

func log(level int, eventOrCallee string, ctx *Context, err error, properties []interface{}) {
	handler := getHandler(level, eventOrCallee, ctx, properties)
	handler.Handle(ctx, err, properties)
}

func getHandler(level int, eventOrCallee string, ctx *Context, properties []interface{}) EventHandler {
	handler, found := handlerCache.Load(eventOrCallee)
	if found {
		return handler.(EventHandler)
	}
	skipFramesCount := 3
	if ctx != nil {
		skipFramesCount = 5
	}
	_, callerFile, callerLine, _ := runtime.Caller(skipFramesCount)
	var handlers EventHandlers
	for _, sink := range EventSinks {
		if !sink.ShouldLog(level, eventOrCallee, properties) {
			continue
		}
		handler := sink.HandlerOf(level, eventOrCallee, callerFile, callerLine, properties)
		handlers = append(handlers, handler)
	}
	switch len(handlers) {
	case 0:
		handler := DefaultEventSink.HandlerOf(level, eventOrCallee, ctx, callerFile, callerLine, properties)
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
