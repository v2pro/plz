package test

import (
	"testing"
	"github.com/v2pro/plz/reflect2"
	"errors"
	"github.com/v2pro/plz/test"
	"github.com/v2pro/plz/countlog"
	"unsafe"
	"github.com/v2pro/plz/test/must"
)

func Test_slice_iface(t *testing.T) {
	t.Run("MakeSlice", testOp(func(api reflect2.API) interface{} {
		valType := api.TypeOf([]error{}).(reflect2.SliceType)
		obj := valType.MakeSlice(5, 10)
		obj.([]error)[0] = errors.New("hello")
		obj.([]error)[4] = errors.New("world")
		return obj
	}))
	t.Run("Set", testOp(func(api reflect2.API) interface{} {
		obj := []error{errors.New("hello"), nil}
		valType := api.TypeOf(obj).(reflect2.SliceType)
		valType.Set(obj, 0, errors.New("hi"))
		valType.Set(obj, 1, errors.New("world"))
		return obj
	}))
	t.Run("UnsafeSet", test.Case(func(ctx *countlog.Context) {
		obj := []error{errors.New("hello"), nil}
		valType := reflect2.TypeOf(obj).(reflect2.SliceType)
		elem0 := errors.New("hi")
		valType.UnsafeSet(reflect2.PtrOf(obj), 0, unsafe.Pointer(&elem0))
		elem1 := errors.New("world")
		valType.UnsafeSet(reflect2.PtrOf(obj), 1, unsafe.Pointer(&elem1))
		must.Equal([]error{elem0, elem1}, obj)
	}))
	t.Run("Get", testOp(func(api reflect2.API) interface{} {
		obj := []error{errors.New("hello"), nil}
		valType := api.TypeOf(obj).(reflect2.SliceType)
		return []interface{}{
			valType.Get(&obj, 0),
			valType.Get(&obj, 1),
			valType.Get(obj, 0),
			valType.Get(obj, 1),
		}
	}))
	t.Run("UnsafeGet", test.Case(func(ctx *countlog.Context) {
		obj := []error{errors.New("hello"), nil}
		valType := reflect2.TypeOf(obj).(reflect2.SliceType)
		elem0 := valType.UnsafeGet(reflect2.PtrOf(obj), 0)
		must.Equal(errors.New("hello"), *(*error)(elem0))
	}))
	t.Run("Append", testOp(func(api reflect2.API) interface{} {
		obj := make([]error, 2, 3)
		obj[0] = errors.New("1")
		obj[1] = errors.New("2")
		valType := api.TypeOf(obj).(reflect2.SliceType)
		obj = valType.Append(obj, errors.New("3")).([]error)
		// will trigger grow
		obj = valType.Append(obj, errors.New("4")).([]error)
		return obj
	}))
	t.Run("UnsafeAppend", test.Case(func(ctx *countlog.Context) {
		obj := make([]error, 2, 3)
		obj[0] = errors.New("1")
		obj[1] = errors.New("2")
		valType := reflect2.TypeOf(obj).(reflect2.SliceType)
		ptr := reflect2.PtrOf(obj)
		elem2 := errors.New("3")
		ptr = valType.UnsafeAppend(ptr, unsafe.Pointer(&elem2))
		elem3 := errors.New("4")
		ptr = valType.UnsafeAppend(ptr, unsafe.Pointer(&elem3))
		must.Equal([]error{
			obj[0], obj[1], elem2, elem3,
		}, valType.PackEFace(ptr))
	}))
}
