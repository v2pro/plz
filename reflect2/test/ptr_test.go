package test

import (
	"testing"
	"github.com/v2pro/plz/reflect2"
)

func Test_ptr(t *testing.T) {
	t.Run("Get from dir", testOp(func(api reflect2.API) interface{} {
		var one = 1
		valType := api.TypeOf(&one).(reflect2.PointerType)
		return valType.Get(&one)
	}))
	t.Run("Get from indir ptr", testOp(func(api reflect2.API) interface{} {
		var one = 1
		var pOne = &one
		valType := api.TypeOf(&pOne).(reflect2.PointerType)
		return valType.Get(&pOne)
	}))
	t.Run("Get from indir map", testOp(func(api reflect2.API) interface{} {
		var m = map[int]int{1:2}
		valType := api.TypeOf(&m).(reflect2.PointerType)
		return valType.Get(&m)
	}))
}
