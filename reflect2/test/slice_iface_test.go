package test

import (
	"testing"
	"github.com/v2pro/plz/reflect2"
	"errors"
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
		valType.Set(&obj, 0, errors.New("hi"))
		valType.Set(&obj, 1, errors.New("world"))
		return obj
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
}