package test

import (
	"testing"
	"github.com/v2pro/plz/reflect2"
)

func Test_slice_array(t *testing.T) {
	t.Run("Set", testOp(func(api reflect2.API) interface{} {
		obj := [][1]int{{}, {}}
		valType := api.TypeOf(obj).(reflect2.SliceType)
		valType.Set(obj, 0, [1]int{1})
		valType.Set(obj, 1, [1]int{2})
		return obj
	}))
	t.Run("Set single ptr struct", testOp(func(api reflect2.API) interface{} {
		obj := [][1]*int{{}, {}}
		valType := api.TypeOf(obj).(reflect2.SliceType)
		valType.Set(obj, 0, [1]*int{})
		valType.Set(obj, 1, [1]*int{})
		return obj
	}))
	t.Run("Set single chan struct", testOp(func(api reflect2.API) interface{} {
		obj := [][1]chan int{{}, {}}
		valType := api.TypeOf(obj).(reflect2.SliceType)
		valType.Set(obj, 0, [1]chan int{})
		valType.Set(obj, 1, [1]chan int{})
		return obj
	}))
	t.Run("Set single func struct", testOp(func(api reflect2.API) interface{} {
		obj := [][1]func(){{}, {}}
		valType := api.TypeOf(obj).(reflect2.SliceType)
		valType.Set(obj, 0, [1]func(){})
		valType.Set(obj, 1, [1]func(){})
		return obj
	}))
}
