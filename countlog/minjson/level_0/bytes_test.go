package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/v2pro/plz/countlog/minjson"
	"reflect"
)

func Test_bytes(t *testing.T) {
	should := require.New(t)
	encoder := minjson.EncoderOf(reflect.TypeOf([]byte(nil)))
	should.Equal(`"hello"`, string(encoder.Encode(nil, ptrOf([]byte("hello")))))
	should.Equal(`"\xe4\xb8\xad\xe6\x96\x87"`, string(encoder.Encode(nil, ptrOf([]byte("中文")))))
	should.Equal(`"\xe4\xb8\xad\n\xe6\x96\x87"`, string(encoder.Encode(nil, ptrOf([]byte("中\n文")))))
}