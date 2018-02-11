package test

import (
	"testing"
	"github.com/v2pro/plz/reflect2"
)

func Test_map(t *testing.T) {
	t.Run("New", testOp(func(api reflect2.API) interface{} {
		valType := api.TypeOf(map[int]int{})
		m := valType.New().(*map[int]int)
		return m
	}))
	t.Run("MakeMap", testOp(func(api reflect2.API) interface{} {
		valType := api.TypeOf(map[int]int{}).(reflect2.MapType)
		m := valType.MakeMap(0).(map[int]int)
		m[2] = 4
		m[3] = 9
		return m
	}))
	t.Run("Set", testOp(func(api reflect2.API) interface{} {
		obj := map[int]int{}
		valType := api.TypeOf(obj).(reflect2.MapType)
		valType.Set(obj, 2, 4)
		valType.Set(obj, 3, 9)
		return obj
	}))
}
