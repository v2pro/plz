package native

import (
	"github.com/stretchr/testify/require"
	"github.com/v2pro/plz"
	"reflect"
	"testing"
)

func Test_struct(t *testing.T) {
	type TestObject struct {
		Field int
	}
	should := require.New(t)
	v := &TestObject{}
	accessor := plz.AccessorOf(reflect.TypeOf(v))
	should.Equal(reflect.Struct, accessor.Kind())
	should.Equal(1, accessor.NumField())
	field := accessor.Field(0)
	should.Equal("Field", field.Name)
	should.Equal(0, field.Accessor.Int(v))
	field.Accessor.SetInt(v, 2)
	should.Equal(2, field.Accessor.Int(v))
}

func Test_struct_tags(t *testing.T) {
	type TestObject struct {
		Field int `json:"field"`
	}
	should := require.New(t)
	v := &TestObject{}
	accessor := plz.AccessorOf(reflect.TypeOf(v))
	should.Equal(map[string]interface{}{"json": "field"}, accessor.Field(0).Tags)
}

func Test_struct_iterate_map(t *testing.T) {
	type TestObject struct {
		Field int
	}
	should := require.New(t)
	v := &TestObject{}
	accessor := plz.AccessorOf(reflect.TypeOf(v))
	should.Equal(reflect.String, accessor.Key().Kind())
	should.Equal(reflect.Interface, accessor.Elem().Kind())
	keys := []string{}
	elems := []int{}
	accessor.IterateMap(v, func(key interface{}, elem interface{}) bool {
		keys = append(keys, accessor.Key().String(key))
		elems = append(elems, accessor.Elem().Int(elem))
		return true
	})
	should.Equal([]string{"Field"}, keys)
	should.Equal([]int{0}, elems)
}
