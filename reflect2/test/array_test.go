package test

import (
	"testing"
	"github.com/v2pro/plz/reflect2"
)

func Test_array(t *testing.T) {
	var pInt = func(val int) *int {
		return &val
	}
	t.Run("New", testOp(func(api reflect2.API) interface{} {
		valType := api.TypeOf([2]int{})
		obj := valType.New()
		(*(obj.(*[2]int)))[0] = 100
		(*(obj.(*[2]int)))[1] = 200
		return obj
	}))
	t.Run("Set", testOp(func(api reflect2.API) interface{} {
		obj := [2]int{}
		valType := api.TypeOf(obj).(reflect2.ArrayType)
		valType.Set(&obj, 0, pInt(100))
		valType.Set(&obj, 1, pInt(200))
		return obj
	}))
	t.Run("Get", testOp(func(api reflect2.API) interface{} {
		obj := [2]int{1, 2}
		valType := api.TypeOf(obj).(reflect2.ArrayType)
		return []interface{} {
			valType.Get(&obj, 0),
			valType.Get(&obj, 1),
		}
	}))
}
