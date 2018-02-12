package test

import (
	"testing"
	"github.com/v2pro/plz/reflect2"
	"github.com/v2pro/plz/test/must"
)

func Test_map_eface_key(t *testing.T) {
	t.Run("Set", testOp(func(api reflect2.API) interface{} {
		obj := map[interface{}]int{}
		valType := api.TypeOf(obj).(reflect2.MapType)
		valType.Set(obj, 2, 4)
		valType.Set(obj, 3, 9)
		valType.Set(obj, nil, 9)
		return obj
	}))
	t.Run("Get", testOp(func(api reflect2.API) interface{} {
		obj := map[interface{}]int{3: 9, 2: 4}
		valType := api.TypeOf(obj).(reflect2.MapType)
		return []interface{}{
			valType.Get(obj, 3),
			valType.Get(obj, 0),
			valType.Get(obj, nil),
			valType.Get(obj, ""),
		}
	}))
	t.Run("Iterate", testOp(func(api reflect2.API) interface{} {
		obj := map[interface{}]int{2: 4, 3: 9}
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
