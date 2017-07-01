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
	should.Equal(acc.Array, accessor.Kind())
	elemAccessor := accessor.Elem()
	elems := []int{}
	accessor.IterateArray(v, func(elem interface{}) bool {
		elems = append(elems, elemAccessor.Int(elem))
		return true
	})
	should.Equal([]int{1, 2, 3}, elems)
}
