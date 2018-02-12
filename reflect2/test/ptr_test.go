package test

import (
	"testing"
	"github.com/v2pro/plz/reflect2"
	"errors"
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
	t.Run("Get from indir eface", testOp(func(api reflect2.API) interface{} {
		var obj interface{} = 123
		valType := api.TypeOf(&obj).(reflect2.PointerType)
		return valType.Get(&obj)
	}))
	t.Run("Get from indir iface", testOp(func(api reflect2.API) interface{} {
		obj := errors.New("hello")
		valType := api.TypeOf(&obj).(reflect2.PointerType)
		return valType.Get(&obj)
	}))
}
