package test

import (
	"testing"
	"github.com/v2pro/plz/reflect2"
	"github.com/v2pro/plz/test"
	"github.com/v2pro/plz/countlog"
	"github.com/v2pro/plz/test/must"
)

func Test_map_elem_bytes(t *testing.T) {
	t.Run("Set", testOp(func(api reflect2.API) interface{} {
		obj := map[int][]byte{}
		valType := api.TypeOf(obj).(reflect2.MapType)
		valType.Set(obj, 2, []byte("hello"))
		valType.Set(obj, 3, nil)
		return obj
	}))
	t.Run("UnsafeSet", test.Case(func(ctx *countlog.Context) {
		obj := map[int][]byte{}
		valType := reflect2.TypeOf(obj).(reflect2.MapType)
		hello := []byte("hello")
		valType.UnsafeSet(reflect2.PtrOf(obj), reflect2.PtrOf(2), reflect2.PtrOf(hello))
		valType.UnsafeSet(reflect2.PtrOf(obj), reflect2.PtrOf(3), nil)
		must.Equal([]byte("hello"), obj[2])
		must.Nil(obj[3])
	}))
	t.Run("UnsafeGet", test.Case(func(ctx *countlog.Context) {
		obj := map[int][]byte{2: []byte("hello")}
		valType := reflect2.TypeOf(obj).(reflect2.MapType)
		elem := valType.UnsafeGet(reflect2.PtrOf(obj), reflect2.PtrOf(2))
		must.Equal([]byte("hello"), valType.Elem().PackEFace(elem))
	}))
}
