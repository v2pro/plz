package countlog

import "github.com/v2pro/plz/countlog/core"

var EventSinks = []EventSink{}

type EventSink interface {
	HandlerOf(level int, eventOrCallee string,
		callerFile string, callerLine int, sample []interface{}) core.EventHandler
	ShouldLog(level int, eventOrCallee string,
		sample []interface{}) bool
}



