package app

import (
	"os"
	"github.com/v2pro/plz/logger"
)

func Run(f func() int, kv ...interface{}) {
	logger.Of("metric", "counter", "begin", "app").
		Info("app begin", kv...)
	defer func() {
		recovered := recover()
		if recovered != nil {
			code := Spi.AfterPanic(recovered, kv)
			Spi.AfterFinish(kv)
			os.Exit(code)
			return
		}
	}()
	code := f()
	Spi.AfterFinish(kv)
	os.Exit(code)
}

var Spi = Config{
	AfterPanic: func(recovered interface{}, kv []interface{}) int {
		logger.Of("metric", "counter", "panic", "app").
			Error("app panic", append(kv, "recovered", recovered)...)
		return 1
	},
	AfterFinish: func(kv []interface{}) {
		logger.Of("metric", "counter", "finish", "app").
			Info("app finish", kv...)
	},
}

type Config struct {
	AfterPanic  func(recovered interface{}, kv []interface{}) int
	AfterFinish func(kv []interface{})
}

func (cfg *Config) Append(newCfg Config) {
	if newCfg.AfterPanic != nil {
		oldAfterPanic := cfg.AfterPanic
		cfg.AfterPanic = func(recovered interface{}, kv []interface{}) int {
			oldAfterPanic(recovered, kv)
			return newCfg.AfterPanic(recovered, kv)
		}
	}
	if newCfg.AfterFinish != nil {
		oldAfterFinish := cfg.AfterFinish
		cfg.AfterFinish = func(kv []interface{}) {
			oldAfterFinish(kv)
			newCfg.AfterFinish(kv)
		}
	}
}
