package test

import (
	"testing"
	"github.com/v2pro/plz/reflect2"
)

func Test_map_elem_bytes(t *testing.T) {
	t.Run("Set", testOp(func(api reflect2.API) interface{} {
		obj := map[int][]byte{}
		valType := api.TypeOf(obj).(reflect2.MapType)
		valType.Set(obj, 2, []byte("hello"))
		valType.Set(obj, 3, nil)
		return obj
	}))
}
