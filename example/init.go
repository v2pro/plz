package example

import "github.com/v2pro/plz/log"

func init() {
	logger.AddLoggerProvider(func(loggerKv []interface{}) logger.Logger {
		return logger.NewStderrLogger(loggerKv, logger.LEVEL_DEBUG)
	})
}
