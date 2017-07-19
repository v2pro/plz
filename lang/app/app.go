package app

import (
	"github.com/v2pro/plz/logging"
	"os"
)

var beginLogger = logging.LoggerOf("metric", "counter", "begin", "app")
var panicLogger = logging.LoggerOf("metric", "counter", "panic", "app")
var finishLogger = logging.LoggerOf("metric", "counter", "finish", "app")

func Run(f func() int, kv ...interface{}) {
	beginLogger.Info("app begin", kv...)
	defer func() {
		recovered := recover()
		if recovered != nil {
			code := -1
			for _, handle := range AfterPanic {
				code = handle(recovered, kv)
			}
			for _, handle := range BeforeFinish {
				handle(kv)
			}
			os.Exit(code)
			return
		}
	}()
	code := f()
	for _, handle := range BeforeFinish {
		handle(kv)
	}
	os.Exit(code)
}

var AfterPanic = []func(recovered interface{}, kv []interface{}) int{
	func(recovered interface{}, kv []interface{}) int {
		panicLogger.Error(nil, "app panic",
			append(kv, "recovered", recovered)...)
		return 1
	},
}

var BeforeFinish = []func(kv []interface{}){
	func(kv []interface{}) {
		finishLogger.Info("app finish", kv...)
	},
}
