package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/v2pro/plz/countlog/minjson"
	"reflect"
	"unsafe"
)

func Test_pointer(t *testing.T) {
	should := require.New(t)
	encoder := minjson.EncoderOf(reflect.TypeOf((*int)(nil)))
	ptr := minjson.PtrOf(1)
	should.Equal("1", string(encoder.Encode(nil, unsafe.Pointer(&ptr))))
	should.Equal("null", string(encoder.Encode(nil, minjson.PtrOf(nil))))
}