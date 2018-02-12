package test

import (
	"testing"
	"github.com/v2pro/plz/reflect2"
	"github.com/v2pro/plz/test/must"
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
	t.Run("Get", testOp(func(api reflect2.API) interface{} {
		obj := map[int]int{3: 9, 2: 4}
		valType := api.TypeOf(obj).(reflect2.MapType)
		return []interface{}{
			valType.Get(obj, 3),
			valType.Get(obj, 0),
		}
	}))
	t.Run("Iterate", testOp(func(api reflect2.API) interface{} {
		obj := map[int]int{2: 4, 3: 9}
		valType := api.TypeOf(obj).(reflect2.MapType)
		iter := valType.Iterate(obj)
		must.Pass(iter.HasNext(), "api", api)
		key1, elem1 := iter.Next()
		must.Pass(iter.HasNext(), "api", api)
		key2, elem2 := iter.Next()
		must.Pass(!iter.HasNext(), "api", api)
		return []interface{}{key1, elem1, key2, elem2}
	}))
}
