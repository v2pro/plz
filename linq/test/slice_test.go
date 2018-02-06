package test

import (
	"testing"
	"github.com/v2pro/plz/test"
	"github.com/v2pro/plz/countlog"
	"github.com/v2pro/plz/linq"
	"github.com/v2pro/plz/test/must"
)

func Test_slice(t *testing.T) {
	t.Run("exact copy", test.Case(func(ctx *countlog.Context) {
		var output []string
		linq.Execute("SELECT word FROM sentence",
			"sentence", []string{"hello", "world"},
			"into", &output)
		must.Equal([]string{"hello"}, output)
	}))
}
