package example

import "github.com/v2pro/plz/log"

func init() {
	log.AddLoggerProvider(func(path []string) log.Logger {
		return log.NewStderrLogger(path, log.LEVEL_DEBUG)
	})
}
