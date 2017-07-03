package native

import (
	"github.com/stretchr/testify/require"
	"github.com/v2pro/plz"
	"reflect"
	"testing"
	"github.com/v2pro/plz/acc"
)

func Test_struct_iterate_array(t *testing.T) {
	type TestObject struct {
		Field1 int
		Field2 int
	}
	should := require.New(t)
	v := &TestObject{1, 2}
	accessor := plz.AccessorOf(reflect.TypeOf(v))
	should.Equal(lang.Struct, accessor.Kind())
	should.Equal(2, accessor.NumField())
	should.Equal("Field1", accessor.Field(0).Name())
	elems := []int{}
	accessor.IterateArray(v, func(index int, elem interface{}) bool {
		elemVal := accessor.Field(index).Accessor().Int(elem)
		elems = append(elems, elemVal)
		return true
	})
	should.Equal([]int{1, 2}, elems)
}

func Test_struct_fill_array(t *testing.T) {
	type TestObject struct {
		Field1 int
		Field2 int
	}
	should := require.New(t)
	v := &TestObject{}
	accessor := plz.AccessorOf(reflect.TypeOf(v))
	accessor.FillArray(v, func(filler lang.ArrayFiller) {
		index, elem := filler.Next()
		should.Equal(0, index)
		accessor.Field(index).Accessor().SetInt(elem, 1)
		filler.Fill()
		index, elem = filler.Next()
		should.Equal(1, index)
		accessor.Field(index).Accessor().SetInt(elem, 2)
		filler.Fill()
	})
	elems := []int{}
	accessor.IterateArray(v, func(index int, elem interface{}) bool {
		elemVal := accessor.Field(index).Accessor().Int(elem)
		elems = append(elems, elemVal)
		return true
	})
	should.Equal([]int{1, 2}, elems)
}

func Test_struct_tags(t *testing.T) {
	type TestObject struct {
		Field int `json:"field"`
	}
	should := require.New(t)
	v := &TestObject{}
	accessor := plz.AccessorOf(reflect.TypeOf(v))
	should.Equal(map[string]interface{}{"json": "field"}, accessor.Field(0).Tags())
}
