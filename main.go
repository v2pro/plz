package plz

import "os"

func RunMain(f func() int) {
	defer func() {
		recovered := recover()
		if recovered != nil {
			code := MainSpi.AfterPanic(recovered)
			MainSpi.AfterFinish()
			os.Exit(code)
			return
		}
	}()
	code := f()
	MainSpi.AfterFinish()
	os.Exit(code)
}

var MainSpi = MainSpiConfig{
	AfterPanic: func(recovered interface{}) int {
		return 1
	},
	AfterFinish: func() {
	},
}

type MainSpiConfig struct {
	AfterPanic  func(recovered interface{}) int
	AfterFinish func()
}

func (cfg *MainSpiConfig) Append(newCfg MainSpiConfig) {
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
