package test

import (
	"testing"
	. "github.com/v2pro/plz/countlog"
	. "github.com/v2pro/plz/check"
	. "github.com/v2pro/plz/check/must"
)

func Test(t *testing.T) {
	t.Run("1 != 2", Case(func(ctx *Context) {
		Check(1 == 2)
	}))
}
