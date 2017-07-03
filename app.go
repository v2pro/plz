package plz

import (
	"github.com/v2pro/plz/lang/app"
)

func RunApp(f func() int) {
	app.Run(f)
}
