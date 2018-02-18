package test

import (
	"testing"
	"github.com/v2pro/plz/reflect2"
	"github.com/v2pro/plz/test"
	"github.com/v2pro/plz/countlog"
	"github.com/v2pro/plz/test/must"
)

func Test_slice(t *testing.T) {
	var pInt = func(val int) *int {
		return &val
	}
	t.Run("New", testOp(func(api reflect2.API) interface{} {
		valType := api.TypeOf([]int{})
		obj := *valType.New().(*[]int)
		obj = append(obj, 1)
		return obj
	}))
	t.Run("IsNil", testOp(func(api reflect2.API) interface{} {
		valType := api.TypeOf([]int{})
		var nilSlice []int
		s := []int{}
		return []interface{}{
			valType.IsNil(&nilSlice),
			valType.IsNil(&s),
			valType.IsNil(nil),
		}
	}))
	t.Run("MakeSlice", testOp(func(api reflect2.API) interface{} {
		valType := api.TypeOf([]int{}).(reflect2.SliceType)
		obj := *(valType.MakeSlice(5, 10).(*[]int))
		obj[0] = 100
		obj[4] = 20
		return obj
	}))
	t.Run("UnsafeMakeSlice", test.Case(func(ctx *countlog.Context) {
		valType := reflect2.TypeOf([]int{}).(reflect2.SliceType)
		obj := valType.UnsafeMakeSlice(5, 10)
		must.Equal(&[]int{0, 0, 0, 0, 0}, valType.PackEFace(obj))
	}))
	t.Run("Set", testOp(func(api reflect2.API) interface{} {
		obj := []int{1, 2}
		valType := api.TypeOf(obj).(reflect2.SliceType)
		valType.Set(&obj, 0, pInt(100))
		valType.Set(&obj, 1, pInt(20))
		return obj
	}))
	t.Run("UnsafeSet", test.Case(func(ctx *countlog.Context) {
		obj := []int{1, 2}
		valType := reflect2.TypeOf(obj).(reflect2.SliceType)
		valType.UnsafeSet(reflect2.PtrOf(obj), 0, reflect2.PtrOf(100))
		valType.UnsafeSet(reflect2.PtrOf(obj), 1, reflect2.PtrOf(10))
		must.Equal([]int{100, 10}, obj)
	}))
	t.Run("Get", testOp(func(api reflect2.API) interface{} {
		obj := []int{1, 2}
		valType := api.TypeOf(obj).(reflect2.SliceType)
		return []interface{}{
			valType.Get(&obj, 1).(*int),
		}
	}))
	t.Run("UnsafeGet", test.Case(func(ctx *countlog.Context) {
		obj := []int{1, 2}
		valType := reflect2.TypeOf(obj).(reflect2.SliceType)
		elem0 := valType.UnsafeGet(reflect2.PtrOf(obj), 0)
		must.Equal(1, *(*int)(elem0))
		elem1 := valType.UnsafeGet(reflect2.PtrOf(obj), 1)
		must.Equal(2, *(*int)(elem1))
	}))
	t.Run("Append", testOp(func(api reflect2.API) interface{} {
		obj := make([]int, 2, 3)
		obj[0] = 1
		obj[1] = 2
		valType := api.TypeOf(obj).(reflect2.SliceType)
		ptr := &obj
		ptr = valType.Append(ptr, pInt(3)).(*[]int)
		// will trigger grow
		ptr = valType.Append(ptr, pInt(4)).(*[]int)
		return ptr
	}))
	t.Run("UnsafeAppend", test.Case(func(ctx *countlog.Context) {
		obj := make([]int, 2, 3)
		obj[0] = 1
		obj[1] = 2
		valType := reflect2.TypeOf(obj).(reflect2.SliceType)
		ptr := reflect2.PtrOf(obj)
		ptr = valType.UnsafeAppend(ptr, reflect2.PtrOf(3))
		ptr = valType.UnsafeAppend(ptr, reflect2.PtrOf(4))
		must.Equal(&[]int{1, 2, 3, 4}, valType.PackEFace(ptr))
	}))
}
