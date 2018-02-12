package test

import (
	"testing"
	"github.com/v2pro/plz/reflect2"
	"github.com/v2pro/plz/test/must"
)

func Test_map_key_eface(t *testing.T) {
	var pEFace = func(val interface{}) interface{} {
		return &val
	}
	var pInt = func(val int) *int {
		return &val
	}
	t.Run("Set", testOp(func(api reflect2.API) interface{} {
		obj := map[interface{}]int{}
		valType := api.TypeOf(obj).(reflect2.MapType)
		valType.Set(&obj, pEFace(2), pInt(4))
		valType.Set(&obj, pEFace(3), pInt(9))
		valType.Set(&obj, pEFace(nil), pInt(9))
		return obj
	}))
	t.Run("Get", testOp(func(api reflect2.API) interface{} {
		obj := map[interface{}]int{3: 9, 2: 4}
		valType := api.TypeOf(obj).(reflect2.MapType)
		return []interface{}{
			valType.Get(&obj, pEFace(3)),
			valType.Get(&obj, pEFace(0)),
			valType.Get(&obj, pEFace(nil)),
			valType.Get(&obj, pEFace("")),
		}
	}))
	t.Run("Iterate", testOp(func(api reflect2.API) interface{} {
		obj := map[interface{}]int{2: 4}
		valType := api.TypeOf(obj).(reflect2.MapType)
		iter := valType.Iterate(obj)
		must.Pass(iter.HasNext(), "api", api)
		key1, elem1 := iter.Next()
		must.Pass(!iter.HasNext(), "api", api)
		return []interface{}{key1, elem1}
	}))
}
