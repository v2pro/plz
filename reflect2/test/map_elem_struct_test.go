package test

import (
	"testing"
	"github.com/v2pro/plz/reflect2"
	"time"
)

func Test_map_elem_struct(t *testing.T) {
	t.Run("Set", testOp(func(api reflect2.API) interface{} {
		obj := map[int]time.Time{}
		valType := api.TypeOf(obj).(reflect2.MapType)
		valType.Set(obj, 2, time.Time{})
		valType.Set(obj, 3, time.Time{})
		return obj
	}))
	t.Run("Set single ptr struct", testOp(func(api reflect2.API) interface{} {
		type TestObject struct {
			Field1 *int
		}
		obj := map[int]TestObject{}
		valType := api.TypeOf(obj).(reflect2.MapType)
		valType.Set(obj, 2, TestObject{})
		valType.Set(obj, 3, TestObject{})
		return obj
	}))
	t.Run("Set single map struct", testOp(func(api reflect2.API) interface{} {
		type TestObject struct {
			Field1 map[int]int
		}
		obj := map[int]TestObject{}
		valType := api.TypeOf(obj).(reflect2.MapType)
		valType.Set(obj, 2, TestObject{})
		valType.Set(obj, 3, TestObject{})
		return obj
	}))
	t.Run("Set single chan struct", testOp(func(api reflect2.API) interface{} {
		type TestObject struct {
			Field1 chan int
		}
		obj := map[int]TestObject{}
		valType := api.TypeOf(obj).(reflect2.MapType)
		valType.Set(obj, 2, TestObject{})
		valType.Set(obj, 3, TestObject{})
		return obj
	}))
	t.Run("Set single func struct", testOp(func(api reflect2.API) interface{} {
		type TestObject struct {
			Field1 func()
		}
		obj := map[int]TestObject{}
		valType := api.TypeOf(obj).(reflect2.MapType)
		valType.Set(obj, 2, TestObject{})
		valType.Set(obj, 3, TestObject{})
		return obj
	}))
}
