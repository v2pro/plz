package example

import "github.com/v2pro/plz/logging"

func init() {
	logging.AddLoggerProvider(func(loggerKv []interface{}) logging.Logger {
		return logging.NewStderrLogger(loggerKv, logging.LEVEL_DEBUG)
	})
}
