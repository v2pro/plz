package native

import (
	"github.com/stretchr/testify/require"
	"github.com/v2pro/plz"
	"reflect"
	"testing"
	"github.com/v2pro/plz/acc"
)

func Test_slice(t *testing.T) {
	should := require.New(t)
	var v interface{} = []int{}
	accessor := plz.AccessorOf(reflect.TypeOf(v))
	should.Equal(lang.Array, accessor.Kind())
	elemAccessor := accessor.Elem()
	elems := []int{}
	accessor.IterateArray(v, func(index int, elem interface{}) bool {
		elems = append(elems, elemAccessor.Int(elem))
		return true
	})
	should.Equal([]int{}, elems)
	// grow one
	accessor.FillArray(v, func(filler lang.ArrayFiller) {
		_, elem := filler.Next()
		accessor.Elem().SetInt(elem, 1)
	})
	elems = []int{}
	// check again
	accessor.IterateArray(v, func(index int, elem interface{}) bool {
		elems = append(elems, elemAccessor.Int(elem))
		return true
	})
	should.Equal([]int{1}, elems)
}

func Test_slice_of_interface(t *testing.T) {
	should := require.New(t)
	v := []interface{}{1, 2, 3}
	accessor := plz.AccessorOf(reflect.TypeOf(v))
	should.Equal(lang.Array, accessor.Kind())
	elemAccessor := accessor.Elem()
	elems := []int{}
	accessor.IterateArray(v, func(index int, elem interface{}) bool {
		elems = append(elems, elemAccessor.Int(elem))
		return true
	})
	should.Equal([]int{1, 2, 3}, elems)
	accessor.FillArray(&v, func(filler lang.ArrayFiller) {
		_, elem := filler.Next()
		elemAccessor.SetInt(elem, 4)
		filler.Fill()
	})
	elems = []int{}
	accessor.IterateArray(v, func(index int, elem interface{}) bool {
		elems = append(elems, elemAccessor.Int(elem))
		return true
	})
	should.Equal([]int{4}, elems)
}
