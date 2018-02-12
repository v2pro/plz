package test

import (
	"testing"
	"github.com/v2pro/plz/reflect2"
)

func Test_map_elem_map(t *testing.T) {
	t.Run("Set", testOp(func(api reflect2.API) interface{} {
		obj := map[int]map[int]int{}
		valType := api.TypeOf(obj).(reflect2.MapType)
		valType.Set(obj, 2, map[int]int{4:4})
		valType.Set(obj, 3, map[int]int{9:9})
		valType.Set(obj, 3, nil)
		return obj
	}))
}