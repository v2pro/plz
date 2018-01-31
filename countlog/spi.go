package countlog

import "unsafe"

var EventSinks = []EventSink{}

type EventSink interface {
	HandlerOf(level int, eventOrCallee string,
		callerFile string, callerLine int, sample []interface{}) EventHandler
	ShouldLog(level int, eventOrCallee string,
		sample []interface{}) bool
}



