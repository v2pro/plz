package native

import (
	"github.com/stretchr/testify/require"
	"github.com/v2pro/plz"
	"reflect"
	"testing"
	"github.com/v2pro/plz/acc"
	"fmt"
)

func Test_map(t *testing.T) {
	should := require.New(t)
	v := map[int]int{}
	accessor := plz.AccessorOf(reflect.TypeOf(v))
	should.Equal(acc.Map, accessor.Kind())
	keys := []interface{}{}
	accessor.IterateMap(v, func(key interface{}, value interface{}) bool {
		keys = append(keys, key)
		return true
	})
	should.Equal([]interface{}{}, keys)
	accessor.FillMap(v, func(filler acc.MapFiller) {
		key, elem := filler.Next()
		accessor.Key().SetInt(key, 1)
		accessor.Elem().SetInt(elem, 2)
		filler.Fill()
	})
	accessor.IterateMap(v, func(key interface{}, value interface{}) bool {
		keys = append(keys, key)
		return true
	})
	should.Equal([]interface{}{1}, keys)
}

func Test_map_reflect(t *testing.T) {
	a := map[string]string{}
	a["hello"] = "world"
	b := reflect.ValueOf(&a).Elem().MapIndex(reflect.ValueOf("hello")).Interface()
	ptr := extractPtrFromEmptyInterface(b)
	*((*string)(ptr)) = "120"
	fmt.Println(a)
}