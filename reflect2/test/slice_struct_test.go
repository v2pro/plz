package test

import (
	"testing"
	"github.com/v2pro/plz/reflect2"
	"time"
)

func Test_slice_struct(t *testing.T) {
	t.Run("Set", testOp(func(api reflect2.API) interface{} {
		obj := []time.Time{{}, {}}
		valType := api.TypeOf(obj).(reflect2.SliceType)
		valType.Set(obj, 0, time.Time{})
		valType.Set(obj, 1, time.Time{})
		return obj
	}))
	t.Run("Set single ptr struct", testOp(func(api reflect2.API) interface{} {
		type TestObject struct {
			Field1 *int
		}
		obj := []TestObject{{}, {}}
		valType := api.TypeOf(obj).(reflect2.SliceType)
		valType.Set(obj, 0, TestObject{})
		valType.Set(obj, 1, TestObject{})
		return obj
	}))
	t.Run("Set single chan struct", testOp(func(api reflect2.API) interface{} {
		type TestObject struct {
			Field1 chan int
		}
		obj := []TestObject{{}, {}}
		valType := api.TypeOf(obj).(reflect2.SliceType)
		valType.Set(obj, 0, TestObject{})
		valType.Set(obj, 1, TestObject{})
		return obj
	}))
	t.Run("Set single func struct", testOp(func(api reflect2.API) interface{} {
		type TestObject struct {
			Field1 func()
		}
		obj := []TestObject{{}, {}}
		valType := api.TypeOf(obj).(reflect2.SliceType)
		valType.Set(obj, 0, TestObject{})
		valType.Set(obj, 1, TestObject{})
		return obj
	}))
}