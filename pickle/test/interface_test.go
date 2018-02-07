package test

import (
	"testing"
	"github.com/v2pro/plz/test"
	"github.com/v2pro/plz/countlog"
	"github.com/v2pro/plz/test/must"
	"github.com/v2pro/plz/pickle"
	"fmt"
)

func Test_interface(t *testing.T) {
	t.Run("slice of struct", test.Case(func(ctx *countlog.Context) {
		type TestObject struct {
			Field int
		}
		elem1 := TestObject{1}
		elem2 := TestObject{2}
		elem3 := TestObject{3}
		obj := []interface{}{elem1, elem2, elem3}
		encoded := must.Call(pickle.Marshal, obj)[0].([]byte)
		decoded := must.Call(pickle.Unmarshal, encoded)[0]
		must.Equal(obj, decoded)
	}))
	t.Run("slice of ptr", test.Case(func(ctx *countlog.Context) {
		type TestObject struct {
			Field int
		}
		elem1 := &TestObject{1}
		elem2 := &TestObject{2}
		elem3 := &TestObject{3}
		obj := []interface{}{elem1, elem2, elem3}
		encoded := must.Call(pickle.Marshal, obj)[0].([]byte)
		fmt.Println(encoded)
		decoded := must.Call(pickle.Unmarshal, encoded)[0]
		must.Equal(obj, decoded)
	}))
}
