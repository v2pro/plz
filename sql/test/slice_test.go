package test

import (
	"testing"
	"github.com/v2pro/plz/test"
	"github.com/v2pro/plz/countlog"
	"github.com/v2pro/plz/sql"
	"github.com/v2pro/plz/test/must"
	"fmt"
)

func Test_slice(t *testing.T) {
	t.Run("exact copy", test.Case(func(ctx *countlog.Context) {
		var output []string
		sql.Q("SELECT word FROM sentence WHERE word != 'hello'",
			"sentence", []string{"hello", "world"},
			"into", &output)
		fmt.Println(output)
		must.Equal([]string{"world"}, output)
	}))
}
