package native

import (
	"github.com/stretchr/testify/require"
	"github.com/v2pro/plz"
	"reflect"
	"testing"
)

func Test_slice_iterate(t *testing.T) {
	should := require.New(t)
	var v interface{} = []int{}
	accessor := plz.AccessorOf(reflect.TypeOf(v))
	should.Equal(reflect.Slice, accessor.Kind())
	elemAccessor := accessor.Elem()
	elems := []int{}
	accessor.IterateArray(v, func(elem interface{}) bool {
		elems = append(elems, elemAccessor.Int(elem))
		return false
	})
	should.Equal([]int{}, elems)
	// grow one
	var tail interface{}
	v, tail = accessor.GrowOne(v, 1)
	elemAccessor.SetInt(tail, 1)
	// check again
	accessor.IterateArray(v, func(elem interface{}) bool {
		elems = append(elems, elemAccessor.Int(elem))
		return false
	})
	should.Equal([]int{1}, elems)
}
