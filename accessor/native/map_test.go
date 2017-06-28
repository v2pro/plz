package native

import (
	"github.com/stretchr/testify/require"
	"github.com/v2pro/plz"
	"reflect"
	"testing"
)

func Test_map(t *testing.T) {
	should := require.New(t)
	v := map[int]int{}
	accessor := plz.AccessorOf(reflect.TypeOf(v))
	should.Equal(reflect.Map, accessor.Kind())
	keys := []interface{}{}
	accessor.IterateMap(v, func(key interface{}, value interface{}) bool {
		keys = append(keys, key)
		return true
	})
	should.Equal([]interface{}{}, keys)
	accessor.SetMap(v, func(key interface{}) {
		accessor.Key().SetInt(key, 1)
	}, func(elem interface{}) {
		accessor.Elem().SetInt(elem, 2)
	})
	accessor.IterateMap(v, func(key interface{}, value interface{}) bool {
		keys = append(keys, key)
		return true
	})
	should.Equal([]interface{}{1}, keys)
}
