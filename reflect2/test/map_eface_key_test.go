package test

import (
	"testing"
	"github.com/v2pro/plz/reflect2"
)

func Test_map_eface_key(t *testing.T) {
	t.Run("Set", testOp(func(api reflect2.API) interface{} {
		obj := map[interface{}]int{}
		valType := api.TypeOf(obj).(reflect2.MapType)
		valType.Set(obj, 2, 4)
		valType.Set(obj, 3, 9)
		return obj
	}))
}