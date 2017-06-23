package app

import "os"

func Run(f func() int) {
	defer func() {
		recovered := recover()
		if recovered != nil {
			code := Spi.AfterPanic(recovered)
			Spi.AfterFinish()
			os.Exit(code)
			return
		}
	}()
	code := f()
	Spi.AfterFinish()
	os.Exit(code)
}

var Spi = Config{
	AfterPanic: func(recovered interface{}) int {
		return 1
	},
	AfterFinish: func() {
	},
}

type Config struct {
	AfterPanic  func(recovered interface{}) int
	AfterFinish func()
}

func (cfg *Config) Append(newCfg Config) {
	if newCfg.AfterPanic != nil {
		oldAfterPanic := cfg.AfterPanic
		cfg.AfterPanic = func(recovered interface{}) int {
			oldAfterPanic(recovered)
			return newCfg.AfterPanic(recovered)
		}
	}
	if newCfg.AfterFinish != nil {
		oldAfterFinish := cfg.AfterFinish
		cfg.AfterFinish = func() {
			oldAfterFinish()
			newCfg.AfterFinish()
		}
	}
}
