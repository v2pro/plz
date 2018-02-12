package test

import (
	"testing"
	"github.com/v2pro/plz/reflect2"
	"github.com/v2pro/plz/test"
	"github.com/v2pro/plz/countlog"
	"github.com/v2pro/plz/test/must"
	"unsafe"
	"github.com/v2pro/plz/test/should"
)

func Test_slice_eface(t *testing.T) {
	t.Run("MakeSlice", testOp(func(api reflect2.API) interface{} {
		valType := api.TypeOf([]interface{}{}).(reflect2.SliceType)
		obj := valType.MakeSlice(5, 10)
		obj.([]interface{})[0] = 100
		obj.([]interface{})[4] = 20
		return obj
	}))
	t.Run("Set", testOp(func(api reflect2.API) interface{} {
		obj := []interface{}{1, nil}
		valType := api.TypeOf(obj).(reflect2.SliceType)
		valType.Set(obj, 0, 100)
		valType.Set(obj, 1, 20)
		return obj
	}))
	t.Run("UnsafeSet", test.Case(func(ctx *countlog.Context) {
		obj := []interface{}{1, 2}
		valType := reflect2.TypeOf(obj).(reflect2.SliceType)
		var elem0 interface{} = 100
		valType.UnsafeSet(reflect2.PtrOf(obj), 0, unsafe.Pointer(&elem0))
		var elem1 interface{} = 10
		valType.UnsafeSet(reflect2.PtrOf(obj), 1, unsafe.Pointer(&elem1))
		must.Equal([]interface{}{100, 10}, obj)
	}))
	t.Run("Get", testOp(func(api reflect2.API) interface{} {
		obj := []interface{}{1, nil}
		valType := api.TypeOf(obj).(reflect2.SliceType)
		return []interface{}{
			valType.Get(&obj, 0),
			valType.Get(&obj, 1),
			valType.Get(obj, 0),
			valType.Get(obj, 1),
		}
	}))
	t.Run("UnsafeGet", test.Case(func(ctx *countlog.Context) {
		obj := []interface{}{1, nil}
		valType := reflect2.TypeOf(obj).(reflect2.SliceType)
		elem0 := valType.UnsafeGet(reflect2.PtrOf(obj), 0)
		must.Equal(1, *(*interface{})(elem0))
	}))
	t.Run("Append", testOp(func(api reflect2.API) interface{} {
		obj := make([]interface{}, 2, 3)
		obj[0] = 1
		obj[1] = 2
		valType := api.TypeOf(obj).(reflect2.SliceType)
		obj = valType.Append(obj, 3).([]interface{})
		// will trigger grow
		obj = valType.Append(obj, 4).([]interface{})
		return obj
	}))
	t.Run("UnsafeAppend", test.Case(func(ctx *countlog.Context) {
		obj := make([]interface{}, 2, 3)
		obj[0] = 1
		obj[1] = 2
		valType := reflect2.TypeOf(obj).(reflect2.SliceType)
		ptr := reflect2.PtrOf(obj)
		var elem2 interface{} = 3
		ptr = valType.UnsafeAppend(ptr, unsafe.Pointer(&elem2))
		var elem3 interface{} = 4
		ptr = valType.UnsafeAppend(ptr, unsafe.Pointer(&elem3))
		should.Equal([]interface{}{1, 2, 3, 4}, valType.PackEFace(ptr))
	}))
}
