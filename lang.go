package plz

import (
	"github.com/v2pro/plz/lang/app"
	"github.com/v2pro/plz/lang/routine"
	"github.com/v2pro/plz/lang/tagging"
)

func RunApp(f func() int) {
	app.Run(f)
}

func Go(oneOff func()) error {
	return routine.Go(oneOff)
}

func GoLongRunning(longRunning func()) error {
	return routine.GoLongRunning(longRunning)
}

func DefineTags(callback interface{}) {
	tagging.Define(callback)
}
