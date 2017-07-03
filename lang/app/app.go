package app

import (
	"github.com/v2pro/plz/logging"
	"os"
)

func Run(f func() int, kv ...interface{}) {
	logging.LoggerOf("metric", "counter", "begin", "app").
		Info("app begin", kv...)
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

type AfterPanicHandler func(recovered interface{}, kv []interface{}) int

var AfterPanic = []AfterPanicHandler{
	func(recovered interface{}, kv []interface{}) int {
		logging.LoggerOf("metric", "counter", "panic", "app").
			Error("app panic", append(kv, "recovered", recovered)...)
		return 1
	},
}

type BeforeFinishHandler func(kv []interface{})

var BeforeFinish = []BeforeFinishHandler{
	func(kv []interface{}) {
		logging.LoggerOf("metric", "counter", "finish", "app").
			Info("app finish", kv...)
	},
}
