package test

import (
	"testing"
	"github.com/v2pro/plz/reflect2"
	"github.com/v2pro/plz/test/must"
	"github.com/v2pro/plz/test"
	"github.com/v2pro/plz/countlog"
	"unsafe"
)

func Test_map_key_ptr(t *testing.T) {
	var pInt = func(val int) *int {
		return &val
	}
	t.Run("Set", testOp(func(api reflect2.API) interface{} {
		obj := map[*int]int{}
		valType := api.TypeOf(obj).(reflect2.MapType)
		valType.Set(obj, pInt(2), 4)
		valType.Set(obj, pInt(2), 9)
		valType.Set(obj, nil, 9)
		return obj[pInt(2)]
	}))
	t.Run("UnsafeSet", test.Case(func(ctx *countlog.Context) {
		obj := map[*int]int{}
		valType := reflect2.TypeOf(obj).(reflect2.MapType)
		v := pInt(2)
		valType.UnsafeSet(reflect2.PtrOf(obj), unsafe.Pointer(v), reflect2.PtrOf(4))
		must.Equal(4, obj[v])
	}))
	t.Run("Get", testOp(func(api reflect2.API) interface{} {
		obj := map[*int]int{pInt(3): 9, pInt(2): 4}
		valType := api.TypeOf(obj).(reflect2.MapType)
		return []interface{}{
			valType.Get(obj, pInt(3)),
			valType.Get(obj, pInt(2)),
			valType.Get(obj, nil),
		}
	}))
	t.Run("Iterate", testOp(func(api reflect2.API) interface{} {
		obj := map[*int]int{pInt(2): 4}
		valType := api.TypeOf(obj).(reflect2.MapType)
		iter := valType.Iterate(obj)
		must.Pass(iter.HasNext(), "api", api)
		key1, elem1 := iter.Next()
		must.Pass(!iter.HasNext(), "api", api)
		return []interface{}{*key1.(*int), elem1}
	}))
}
