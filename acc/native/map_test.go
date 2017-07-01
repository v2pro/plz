package native

import (
	"github.com/stretchr/testify/require"
	"github.com/v2pro/plz"
	"reflect"
	"testing"
	"github.com/v2pro/plz/acc"
)

func Test_map(t *testing.T) {
	should := require.New(t)
	v := map[int]int{}
	accessor := plz.AccessorOf(reflect.TypeOf(v))
	should.Equal(acc.Map, accessor.Kind())
	keys := []int{}
	elems := []int{}
	accessor.IterateMap(v, func(key interface{}, elem interface{}) bool {
		keys = append(keys, accessor.Key().Int(key))
		elems = append(elems, accessor.Elem().Int(elem))
		return true
	})
	should.Equal([]int{}, keys)
	should.Equal([]int{}, elems)
	accessor.FillMap(v, func(filler acc.MapFiller) {
		key, elem := filler.Next()
		accessor.Key().SetInt(key, 1)
		accessor.Elem().SetInt(elem, 2)
		filler.Fill()
	})
	accessor.IterateMap(v, func(key interface{}, elem interface{}) bool {
		keys = append(keys, accessor.Key().Int(key))
		elems = append(elems, accessor.Elem().Int(elem))
		return true
	})
	should.Equal([]int{1}, keys)
	should.Equal([]int{2}, elems)
}

func Test_map_of_interface(t *testing.T) {
	should := require.New(t)
	v := map[string]interface{}{
		"hello": "world",
	}
	accessor := plz.AccessorOf(reflect.TypeOf(v))
	should.Equal(acc.Interface, accessor.Elem().Kind())
	keys := []string{}
	elems := []string{}
	accessor.IterateMap(v, func(key interface{}, elem interface{}) bool {
		keys = append(keys, accessor.Key().String(key))
		elems = append(elems, accessor.Elem().String(elem))
		return true
	})
	should.Equal([]string{"hello"}, keys)
}
