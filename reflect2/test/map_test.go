package test

import (
	"testing"
	"github.com/v2pro/plz/reflect2"
	"github.com/v2pro/plz/test/must"
	"github.com/v2pro/plz/countlog"
	"github.com/v2pro/plz/test"
	"unsafe"
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
	t.Run("UnsafeMakeMap", test.Case(func(ctx *countlog.Context) {
		valType := reflect2.TypeOf(map[int]int{}).(reflect2.MapType)
		m := *(*map[int]int)(valType.UnsafeMakeMap(0))
		m[2] = 4
		m[3] = 9
	}))
	t.Run("PackEFace", test.Case(func(ctx *countlog.Context) {
		valType := reflect2.TypeOf(map[int]int{}).(reflect2.MapType)
		m := valType.UnsafeMakeMap(0)
		must.Equal(map[int]int{}, valType.PackEFace(unsafe.Pointer(m)))
	}))
	t.Run("Set", testOp(func(api reflect2.API) interface{} {
		obj := map[int]int{}
		valType := api.TypeOf(obj).(reflect2.MapType)
		valType.Set(obj, 2, 4)
		valType.Set(obj, 3, 9)
		return obj
	}))
	t.Run("UnsafeSet", test.Case(func(ctx *countlog.Context) {
		obj := map[int]int{}
		valType := reflect2.TypeOf(obj).(reflect2.MapType)
		valType.UnsafeSet(unsafe.Pointer(&obj), reflect2.PtrOf(2), reflect2.PtrOf(4))
		must.Equal(map[int]int{2: 4}, obj)
	}))
	t.Run("Get", testOp(func(api reflect2.API) interface{} {
		obj := map[int]int{3: 9, 2: 4}
		valType := api.TypeOf(obj).(reflect2.MapType)
		return []interface{}{
			valType.Get(obj, 3),
			valType.Get(obj, 0),
		}
	}))
	t.Run("UnsafeGet", test.Case(func(ctx *countlog.Context) {
		obj := map[int]int{3: 9, 2: 4}
		valType := reflect2.TypeOf(obj).(reflect2.MapType)
		elem := valType.UnsafeGet(unsafe.Pointer(&obj), reflect2.PtrOf(3))
		must.Equal(9, *(*int)(elem))
	}))
	t.Run("Iterate", testOp(func(api reflect2.API) interface{} {
		obj := map[int]int{2: 4}
		valType := api.TypeOf(obj).(reflect2.MapType)
		iter := valType.Iterate(obj)
		must.Pass(iter.HasNext(), "api", api)
		key1, elem1 := iter.Next()
		must.Pass(!iter.HasNext(), "api", api)
		return []interface{}{key1, elem1}
	}))
	t.Run("UnsafeIterate", test.Case(func(ctx *countlog.Context) {
		obj := map[int]int{2: 4}
		valType := reflect2.TypeOf(obj).(reflect2.MapType)
		iter := valType.UnsafeIterate(unsafe.Pointer(&obj))
		must.Pass(iter.HasNext())
		key, elem := iter.UnsafeNext()
		must.Equal(2, *(*int)(key))
		must.Equal(4, *(*int)(elem))
	}))
}
