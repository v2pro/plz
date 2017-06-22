package plz

type Routine struct {
	OneOff func()
}

func (cfg Routine) Go() {
	go func() {
		defer func() {
			recovered := recover()
			for _, handlePanic := range panicHandlers {
				handlePanic(&cfg, recovered)
			}
		}()
		cfg.OneOff()
	}()
}

func Go(oneOff func()) {
	Routine{OneOff: oneOff}.Go()
}

type HandlePanic func(routine *Routine, recovered interface{})

var panicHandlers = []HandlePanic{}

func RegisterPanicHandler(handlePanic HandlePanic) {
	panicHandlers = append(panicHandlers, handlePanic)
}
