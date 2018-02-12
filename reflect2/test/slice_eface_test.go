package test

import (
	"testing"
	"github.com/v2pro/plz/reflect2"
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
		valType.Set(&obj, 0, 100)
		valType.Set(&obj, 1, 20)
		return obj
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
}
