package native

import (
	"testing"
	"reflect"
	"github.com/json-iterator/go/require"
	"github.com/v2pro/plz"
)

func Test_int(t *testing.T) {
	should := require.New(t)
	v := int(0)
	accessor := plz.AccessorOf(&v)
	should.Equal(reflect.Int, accessor.Kind())
	should.Equal(0, accessor.Int(&v))
	accessor.SetInt(&v, 2)
	should.Equal(2, accessor.Int(&v))
}
