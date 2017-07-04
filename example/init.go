package example

import "github.com/v2pro/plz/logging"

func init() {
	logging.Providers = append(logging.Providers, func(loggerKv []interface{}) logging.Logger {
		return logging.NewStderrLogger(loggerKv, logging.LEVEL_DEBUG)
	})
}
