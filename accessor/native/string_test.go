package native

import (
	"testing"
	"github.com/v2pro/plz"
	"reflect"
	"github.com/stretchr/testify/require"
)

func Test_string(t *testing.T) {
	should := require.New(t)
	directV := string("")
	v := &directV
	accessor := plz.AccessorOf(reflect.TypeOf(v))
	should.Equal(reflect.Int, accessor.Kind())
	should.Equal("", accessor.String(v))
	accessor.SetString(v, "hello")
	should.Equal("hello", accessor.String(v))
}

