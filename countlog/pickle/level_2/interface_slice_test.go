package test

import (
	"testing"
	"github.com/v2pro/plz/test"
	"github.com/v2pro/plz/countlog"
	"github.com/v2pro/plz/countlog/pickle"
	"github.com/v2pro/plz/test/must"
	"fmt"
)

func Test_interface_slice(t *testing.T) {
	t.Run("encode decode", test.Case(func(ctx *countlog.Context) {
		output := must.Call(pickle.Marshal, []interface{}{1, 2, 3})[0].([]byte)
		fmt.Println(output)
	}))
}
