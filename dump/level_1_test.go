package dump_test

import (
	"testing"
	"github.com/v2pro/plz/test"
	"github.com/v2pro/plz/countlog"
	"github.com/v2pro/plz/test/must"
	"github.com/v2pro/plz/dump"
)

func Test_level1(t *testing.T) {
	t.Run("string", test.Case(func(ctx *countlog.Context) {
		must.JsonEqual(`{
		"__root__": {
			"type": "string",
			"data": {
				"__ptr__": "{ptr1}"
			}
		},
		"{ptr1}": {
			"data": {
				"__ptr__": "{ptr2}"
			},
			"len": 5
		},
		"{ptr2}": "hello"}`, dump.Var{"hello"}.String())
	}))
}