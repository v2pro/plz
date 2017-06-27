package plz

import (
	"github.com/v2pro/plz/logger"
)

func LoggerOf(loggerKv ...interface{}) logger.Logger {
	return logger.Of(loggerKv...)
}
