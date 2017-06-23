package plz

import (
	"github.com/v2pro/plz/app"
)

func RunApp(f func() int) {
	app.Run(f)
}
