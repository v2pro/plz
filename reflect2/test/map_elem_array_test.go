package test

import (
	"testing"
	"github.com/v2pro/plz/reflect2"
)

func Test_map_elem_array(t *testing.T) {
	t.Run("Set", testOp(func(api reflect2.API) interface{} {
		obj := map[int][2]*int{}
		valType := api.TypeOf(obj).(reflect2.MapType)
		valType.Set(obj, 2, [2]*int{(*int)(reflect2.PtrOf(1)), (*int)(reflect2.PtrOf(2))})
		valType.Set(obj, 3, [2]*int{(*int)(reflect2.PtrOf(3)), (*int)(reflect2.PtrOf(4))})
		return obj
	}))
	t.Run("Set zero length array", testOp(func(api reflect2.API) interface{} {
		obj := map[int][0]*int{}
		valType := api.TypeOf(obj).(reflect2.MapType)
		valType.Set(obj, 2, [0]*int{})
		valType.Set(obj, 3, [0]*int{})
		return obj
	}))
	t.Run("Set single ptr array", testOp(func(api reflect2.API) interface{} {
		obj := map[int][1]*int{}
		valType := api.TypeOf(obj).(reflect2.MapType)
		valType.Set(obj, 2, [1]*int{(*int)(reflect2.PtrOf(1))})
		valType.Set(obj, 3, [1]*int{})
		return obj
	}))
	t.Run("Set single chan array", testOp(func(api reflect2.API) interface{} {
		obj := map[int][1]chan int{}
		valType := api.TypeOf(obj).(reflect2.MapType)
		valType.Set(obj, 2, [1]chan int{})
		valType.Set(obj, 3, [1]chan int{})
		return obj
	}))
	t.Run("Set single func array", testOp(func(api reflect2.API) interface{} {
		obj := map[int][1]func(){}
		valType := api.TypeOf(obj).(reflect2.MapType)
		valType.Set(obj, 2, [1]func(){})
		valType.Set(obj, 3, [1]func(){})
		return obj
	}))
}
