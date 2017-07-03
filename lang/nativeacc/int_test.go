package nativeacc

import (
	"github.com/json-iterator/go/require"
	"github.com/v2pro/plz"
	"reflect"
	"testing"
	"github.com/v2pro/plz/lang"
)

func Test_int(t *testing.T) {
	should := require.New(t)
	directV := int(1)
	v := &directV
	accessor := plz.AccessorOf(reflect.TypeOf(v))
	should.Equal(lang.Int, accessor.Kind())
	should.Equal(1, accessor.Int(v))
	accessor.SetInt(v, 2)
	should.Equal(2, accessor.Int(v))

	accessor = plz.AccessorOf(reflect.TypeOf(directV))
	should.Equal(lang.Int, accessor.Kind())
	should.Equal(2, accessor.Int(directV))
	should.Panics(func() {
		accessor.SetInt(directV, 3)
	})
}
