package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/v2pro/plz/countlog/minjson"
	"reflect"
)

func Test_pointer(t *testing.T) {
	should := require.New(t)
	encoder := minjson.EncoderOf(reflect.TypeOf((*int)(nil)))
	should.Equal("1", string(encoder.Encode(nil, minjson.PtrOf(1))))
	should.Equal("null", string(encoder.Encode(nil, minjson.PtrOf(nil))))
}