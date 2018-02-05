package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"reflect"
	"github.com/v2pro/plz/msgfmt/njson"
)

func Test_bytes(t *testing.T) {
	should := require.New(t)
	encoder := njson.EncoderOf(reflect.TypeOf([]byte(nil)))
	should.Equal(`"hello"`, string(encoder.Encode(nil, njson.PtrOf([]byte("hello")))))
	should.Equal(`"\xe4\xb8\xad\xe6\x96\x87"`, string(encoder.Encode(nil, njson.PtrOf([]byte("中文")))))
	should.Equal(`"\xe4\xb8\xad\n\xe6\x96\x87"`, string(encoder.Encode(nil, njson.PtrOf([]byte("中\n文")))))
}