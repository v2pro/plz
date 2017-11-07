// +build koala_go

package witch

import (
	"runtime"
	"github.com/v2pro/plz/countlog"
)

func setCurrentGoRoutineIsKoala() {
	runtime.SetCurrentGoRoutineIsKoala()
}