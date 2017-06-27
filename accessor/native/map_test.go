package native

import (
	"github.com/stretchr/testify/require"
	"github.com/v2pro/plz"
	"reflect"
	"testing"
)

func Test_map_iterate(t *testing.T) {
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
	accessor.SetMapIndex(v, 1, 2)
	accessor.IterateMap(v, func(key interface{}, value interface{}) bool {
		keys = append(keys, key)
		return true
	})
	should.Equal([]interface{}{1}, keys)
}

func Test_map_value_accessor(t *testing.T) {
	should := require.New(t)
	v := map[int]int{}
	accessor := plz.AccessorOf(reflect.TypeOf(v))
	should.Equal(2, accessor.Key().Int(2))
	should.Equal(2, accessor.Elem().Int(2))
}
