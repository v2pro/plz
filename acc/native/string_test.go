package native

import (
	"testing"
	"github.com/v2pro/plz"
	"reflect"
	"github.com/stretchr/testify/require"
	"github.com/v2pro/plz/acc"
)

func Test_string(t *testing.T) {
	should := require.New(t)
	directV := string("")
	v := &directV
	accessor := plz.AccessorOf(reflect.TypeOf(v))
	should.Equal(acc.String, accessor.Kind())
	should.Equal("", accessor.String(v))
	accessor.SetString(v, "hello")
	should.Equal("hello", accessor.String(v))

	accessor = plz.AccessorOf(reflect.TypeOf(directV))
	should.Equal(acc.String, accessor.Kind())
	should.Equal("hello", accessor.String(directV))
	should.Panics(func() {
		accessor.SetString(directV, "hello")
	})
}

