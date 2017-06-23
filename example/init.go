package example

import "github.com/v2pro/plz/log"

func init() {
	log.AddLoggerProvider(func(loggerKv []interface{}) log.Logger {
		return log.NewStderrLogger(loggerKv, log.LEVEL_DEBUG)
	})
}
