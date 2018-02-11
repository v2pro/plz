package test

import (
	"testing"
	"github.com/v2pro/plz/reflect2"
)

func Test_array(t *testing.T) {
	t.Run("New", testOp(func(api reflect2.API) interface{} {
		valType := api.TypeOf([2]int{})
		obj := valType.New()
		(*(obj.(*[2]int)))[0] = 100
		(*(obj.(*[2]int)))[1] = 200
		return obj
	}))
	t.Run("Set", testOp(func(api reflect2.API) interface{} {
		obj := [2]int{}
		valType := api.TypeOf(obj).(reflect2.ListType)
		valType.Set(&obj, 0, 100)
		valType.Set(&obj, 1, 200)
		return obj
	}))
}
