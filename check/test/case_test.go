package test

import (
	"testing"
	"github.com/v2pro/plz/countlog"
	"github.com/v2pro/plz/check"
	"github.com/v2pro/plz/check/should"
	"github.com/v2pro/plz/check/must"
)

func Test_case(t *testing.T) {
	t.Run("should.Pass will not exit", check.Case(func(ctx *countlog.Context) {
		should.Pass(1 == 1)
		should.Pass(1 == 2)
		ctx.Info("hello")
	}))
	t.Run("must.Pass will exit", check.Case(func(ctx *countlog.Context) {
		must.Pass(1 == 1)
		must.Pass(1 == 2)
		ctx.Info("hello")
	}))
	t.Run("multiline", check.Case(func(ctx *countlog.Context) {
		var f = func(i int) int {return i}
		must.Pass(struct{ i int }{
			f(100),
		} == struct{ i int }{
			f(101),
		})
	}))
}
