package native

import (
	"testing"
	"github.com/stretchr/testify/require"
	"reflect"
	"github.com/v2pro/plz/acc"
)

func Test_empty_interface_of_single_value(t *testing.T) {
	should := require.New(t)
	var directV interface{}
	accessor := lang.AccessorOf(reflect.TypeOf(&directV))
	should.Equal(lang.Interface, accessor.Kind())
	accessor.SetInt(&directV, 10)
	should.Equal(10, directV)
	directV = "hello"
	should.Panics(func() {
		accessor.Int(&directV)
	})
	should.Equal("hello", accessor.String(&directV))
	accessor.SetString(&directV, "world")
	should.Equal("world", accessor.String(&directV))
	accessor.SetInt(&directV, 20)
	should.Equal(20, directV)
	should.Equal(20, accessor.Int(&directV))
}

func Test_empty_interface_not_nil_fill_array(t *testing.T) {
	should := require.New(t)
	var directV interface{} = []int{}
	accessor := lang.AccessorOf(reflect.TypeOf(&directV))
	v, vAccessor := accessor.PtrElem(&directV)
	vAccessor.FillArray(v, func(filler lang.ArrayFiller) {
		_, elem := filler.Next()
		vAccessor.Elem().SetInt(elem, 1)
		filler.Fill()
	})
	should.Equal([]int{1}, directV)
}

func Test_empty_interface_nil_fill_array(t *testing.T) {
	should := require.New(t)
	var directV interface{}
	accessor := lang.AccessorOf(reflect.TypeOf(&directV))
	v, vAccessor := accessor.PtrElem(&directV)
	should.Nil(v)
	should.Nil(vAccessor)
	accessor.SetPtrElem(&directV, []int{2, 5, 0})
	v, vAccessor = accessor.PtrElem(&directV)
	should.Equal(reflect.Slice, reflect.TypeOf(v).Kind())
	should.NotNil(vAccessor)
	vAccessor.FillArray(v, func(filler lang.ArrayFiller) {
		_, elem := filler.Next()
		vAccessor.Elem().SetInt(elem, 1)
		filler.Fill()
	})
	should.Equal([]int{1}, directV)
}
