package plz

import (
	"github.com/v2pro/plz/logging"
)

func LoggerOf(loggerKv ...interface{}) logging.Logger {
	return logging.LoggerOf(loggerKv...)
}
