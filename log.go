package plz

import "github.com/v2pro/plz/log"

func Logger(path ...string) log.Logger {
	return log.GetLogger(path...)
}
