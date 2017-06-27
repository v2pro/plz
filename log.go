package plz

import "github.com/v2pro/plz/log"

func Logger(loggerKv ...interface{}) logger.Logger {
	return logger.GetLogger(loggerKv...)
}
