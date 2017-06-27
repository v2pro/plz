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
