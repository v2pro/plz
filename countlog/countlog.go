package countlog

import (
	"sync"
	"unsafe"
	"runtime"
	"time"
	"github.com/v2pro/plz/countlog/spi"
)

// push event out

const LevelTrace = spi.LevelTrace
const LevelDebug = spi.LevelDebug
const LevelInfo = spi.LevelInfo
const LevelWarn = spi.LevelWarn
const LevelError = spi.LevelError
const LevelFatal = spi.LevelFatal

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

func SetMinLevel(level int) {
	spi.MinLevel = level
	spi.SuccinctLevel = level + 10
}

func ShouldLog(level int) bool {
	return level >= spi.MinLevel
}

func Trace(event string, properties ...interface{}) {
	if LevelTrace < spi.MinLevel {
		return
	}
	log(LevelTrace, event, nil, nil, properties)
}

func TraceCall(event string, err error, properties ...interface{}) {
	if err != nil {
		log(LevelError, event, nil, err, properties)
		return
	}
	if LevelTrace < spi.MinLevel {
		return
	}
	log(LevelTrace, event, nil, err, properties)
}

func Debug(event string, properties ...interface{}) {
	if LevelDebug < spi.MinLevel {
		return
	}
	log(LevelDebug, event, nil, nil, properties)
}

func DebugCall(event string, err error, properties ...interface{}) {
	if err != nil {
		log(LevelError, event, nil, err, properties)
		return
	}
	if LevelDebug < spi.MinLevel {
		return
	}
	log(LevelDebug, event, nil, err, properties)
}

func Info(event string, properties ...interface{}) {
	if LevelInfo < spi.MinLevel {
		return
	}
	log(LevelInfo, event, nil, nil, properties)
}

func InfoCall(event string, err error, properties ...interface{}) {
	if err != nil {
		log(LevelError, event, nil, err, properties)
		return
	}
	if LevelInfo < spi.MinLevel {
		return
	}
	log(LevelInfo, event, nil, err, properties)
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
	if recovered == nil {
		return nil
	}
	buf := make([]byte, 1<<16)
	runtime.Stack(buf, false)
	if len(properties) > 0 {
		properties = append(properties, "err", recovered, "stacktrace", string(buf))
		Fatal("event!panic", properties...)
	} else {
		Fatal("event!panic", "err", recovered, "stacktrace", string(buf))
	}
	return recovered
}

var handlerCache = &sync.Map{}

func log(level int, eventName string, ctx *Context, err error, properties []interface{}) {
	handler := getHandler(level, eventName, ctx, properties)
	event := &spi.Event{
		Level:      level,
		Context:    ctx,
		Error:      err,
		Timestamp:  time.Now(),
		Properties: properties,
	}
	ptr := unsafe.Pointer(event)
	handler.Handle(castEvent(uintptr(ptr)))
}

func castEmptyInterfaces(ptr uintptr) []interface{} {
	return *(*[]interface{})(unsafe.Pointer(ptr))
}

func castEvent(ptr uintptr) *spi.Event {
	return (*spi.Event)(unsafe.Pointer(ptr))
}
func castString(ptr uintptr) string {
	return *(*string)(unsafe.Pointer(ptr))
}

func getHandler(level int, eventOrCallee string, ctx *Context, properties []interface{}) spi.EventHandler {
	handler, found := handlerCache.Load(eventOrCallee)
	if found {
		return handler.(spi.EventHandler)
	}
	return newHandler(level, eventOrCallee, ctx, properties)
}

func newHandler(level int, eventOrCalleeObj string, ctx *Context, properties []interface{}) spi.EventHandler {
	ptr := unsafe.Pointer(&eventOrCalleeObj)
	eventOrCallee := castString(uintptr(ptr))
	skipFramesCount := 3
	if ctx != nil {
		skipFramesCount = 5
	}
	_, callerFile, callerLine, _ := runtime.Caller(skipFramesCount)
	site := &spi.LogSite{
		Level:  level,
		Event:  eventOrCallee,
		File:   callerFile,
		Line:   callerLine,
		Sample: properties,
	}
	var handlers spi.EventHandlers
	for _, sink := range EventSinks {
		handler := sink.HandlerOf(site)
		if handler == nil {
			continue
		}
		handlers = append(handlers, handler)
	}
	switch len(handlers) {
	case 0:
		handler := DevelopmentEventSink.HandlerOf(site)
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

// TODO: remove StateExporter in favor of expvar

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
