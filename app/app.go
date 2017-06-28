package app

import (
	"github.com/v2pro/plz/logging"
	"os"
)

func Run(f func() int, kv ...interface{}) {
	logging.LoggerOf("metric", "counter", "begin", "app").
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
		logging.LoggerOf("metric", "counter", "panic", "app").
			Error("app panic", append(kv, "recovered", recovered)...)
		return 1
	},
	AfterFinish: func(kv []interface{}) {
		logging.LoggerOf("metric", "counter", "finish", "app").
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
