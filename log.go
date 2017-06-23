package plz

import "github.com/v2pro/plz/log"

func Logger(loggerKv ...interface{}) log.Logger {
	return log.GetLogger(loggerKv...)
}
