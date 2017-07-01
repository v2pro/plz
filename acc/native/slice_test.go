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
	should.Equal(acc.Array, accessor.Kind())
	elemAccessor := accessor.Elem()
	elems := []int{}
	accessor.IterateArray(v, func(elem interface{}) bool {
		elems = append(elems, elemAccessor.Int(elem))
		return true
	})
	should.Equal([]int{}, elems)
	// grow one
	accessor.FillArray(v, func(filler acc.ArrayFiller) {
		accessor.Elem().SetInt(filler.Next(), 1)
	})
	// check again
	accessor.IterateArray(v, func(elem interface{}) bool {
		elems = append(elems, elemAccessor.Int(elem))
		return false
	})
	should.Equal([]int{1}, elems)
}

func Test_slice_of_interface(t *testing.T) {
	should := require.New(t)
	v := []interface{}{1, 2, 3}
	accessor := plz.AccessorOf(reflect.TypeOf(v))
	should.Equal(acc.Array, accessor.Kind())
	elemAccessor := accessor.Elem()
	elems := []int{}
	accessor.IterateArray(v, func(elem interface{}) bool {
		elems = append(elems, elemAccessor.Int(elem))
		return false
	})
}
