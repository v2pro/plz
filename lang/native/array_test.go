package native

import (
	"testing"
	"github.com/v2pro/plz"
	"reflect"
	"github.com/v2pro/plz/acc"
	"github.com/stretchr/testify/require"
)

func Test_array(t *testing.T) {
	should := require.New(t)
	var v interface{} = [3]int{1, 2, 3}
	accessor := plz.AccessorOf(reflect.TypeOf(v))
	should.Equal(lang.Array, accessor.Kind())
	elemAccessor := accessor.Elem()
	elems := []int{}
	accessor.IterateArray(v, func(index int, elem interface{}) bool {
		elems = append(elems, elemAccessor.Int(elem))
		return true
	})
	should.Equal([]int{1, 2, 3}, elems)
}

func Test_array_append(t *testing.T) {
	should := require.New(t)
	v := [3]int{1, 2, 3}
	accessor := plz.AccessorOf(reflect.TypeOf(v))
	should.Equal(lang.Array, accessor.Kind())
	elemAccessor := accessor.Elem()
	accessor.FillArray(&v, func(filler lang.ArrayFiller) {
		_, elem := filler.Next()
		should.NotNil(elem)
		elemAccessor.SetInt(elem, 3)
		_, elem = filler.Next()
		should.NotNil(elem)
		elemAccessor.SetInt(elem, 2)
		_, elem = filler.Next()
		should.NotNil(elem)
		elemAccessor.SetInt(elem, 1)
		_, elem = filler.Next()
		should.Nil(elem)
	})
	should.Equal([]int{3, 2, 1}, v[:])
}
