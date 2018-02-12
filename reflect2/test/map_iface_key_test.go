package test

import (
	"testing"
	"github.com/v2pro/plz/reflect2"
	"github.com/v2pro/plz/test/must"
)

type intError int

func (err intError) Error() string {
	return ""
}

func Test_map_iface_key(t *testing.T) {
	t.Run("Set", testOp(func(api reflect2.API) interface{} {
		obj := map[error]int{}
		valType := api.TypeOf(obj).(reflect2.MapType)
		valType.Set(obj, intError(2), 4)
		valType.Set(obj, intError(2), 9)
		valType.Set(obj, nil, 9)
		must.Panic(func() {
			valType.Set(obj, "", 9)
		})
		return obj
	}))
	t.Run("Get", testOp(func(api reflect2.API) interface{} {
		obj := map[error]int{intError(3): 9, intError(2): 4}
		valType := api.TypeOf(obj).(reflect2.MapType)
		must.Panic(func() {
			valType.Get(obj, "")
		})
		return []interface{}{
			valType.Get(obj, intError(3)),
			valType.Get(obj, nil),
		}
	}))
}
