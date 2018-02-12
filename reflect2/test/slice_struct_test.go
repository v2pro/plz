package test

import (
	"testing"
	"github.com/v2pro/plz/reflect2"
)

func Test_slice_struct(t *testing.T) {
	var pInt = func(val int) *int {
		return &val
	}
	t.Run("Set", testOp(func(api reflect2.API) interface{} {
		type TestObject struct {
			Field1 float64
			Field2 float64
		}
		obj := []TestObject{{}, {}}
		valType := api.TypeOf(obj).(reflect2.SliceType)
		valType.Set(obj, 0, &TestObject{1, 3})
		valType.Set(obj, 1, &TestObject{2, 4})
		return obj
	}))
	t.Run("Set single ptr struct", testOp(func(api reflect2.API) interface{} {
		type TestObject struct {
			Field1 *int
		}
		obj := []TestObject{{}, {}}
		valType := api.TypeOf(obj).(reflect2.SliceType)
		valType.Set(obj, 0, &TestObject{pInt(1)})
		valType.Set(obj, 1, &TestObject{pInt(2)})
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