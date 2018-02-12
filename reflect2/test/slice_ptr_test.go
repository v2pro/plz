package test

import (
	"testing"
	"github.com/v2pro/plz/reflect2"
	"github.com/v2pro/plz/test"
	"github.com/v2pro/plz/countlog"
	"unsafe"
	"github.com/v2pro/plz/test/must"
)

func Test_slice_ptr(t *testing.T) {
	var pInt = func(val int) *int {
		return &val
	}
	t.Run("MakeSlice", testOp(func(api reflect2.API) interface{} {
		valType := api.TypeOf([]*int{}).(reflect2.SliceType)
		obj := valType.MakeSlice(5, 10)
		obj.([]*int)[0] = pInt(1)
		obj.([]*int)[4] = pInt(5)
		return obj
	}))
	t.Run("Set", testOp(func(api reflect2.API) interface{} {
		obj := []*int{pInt(1), nil}
		valType := api.TypeOf(obj).(reflect2.SliceType)
		valType.Set(obj, 0, pInt(2))
		valType.Set(obj, 1, pInt(3))
		return obj
	}))
	t.Run("UnsafeSet", test.Case(func(ctx *countlog.Context) {
		obj := []*int{pInt(1), nil}
		valType := reflect2.TypeOf(obj).(reflect2.SliceType)
		valType.UnsafeSet(reflect2.PtrOf(obj), 0, unsafe.Pointer(pInt(2)))
		valType.UnsafeSet(reflect2.PtrOf(obj), 1, unsafe.Pointer(pInt(1)))
		must.Equal([]*int{pInt(2), pInt(1)}, obj)
	}))
	t.Run("Get", testOp(func(api reflect2.API) interface{} {
		obj := []*int{pInt(1), nil}
		valType := api.TypeOf(obj).(reflect2.SliceType)
		return []interface{}{
			valType.Get(&obj, 0),
			valType.Get(&obj, 1),
			valType.Get(obj, 0),
			valType.Get(obj, 1),
		}
	}))
	t.Run("Append", testOp(func(api reflect2.API) interface{} {
		obj := make([]*int, 2, 3)
		obj[0] = pInt(1)
		obj[1] = pInt(2)
		valType := api.TypeOf(obj).(reflect2.SliceType)
		obj = valType.Append(obj, pInt(3)).([]*int)
		// will trigger grow
		obj = valType.Append(obj, pInt(4)).([]*int)
		return obj
	}))
}
