package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/v2pro/plz/countlog/minjson"
	"reflect"
)

func Test_string(t *testing.T) {
	should := require.New(t)
	encoder := minjson.EncoderOf(reflect.TypeOf(""))
	should.Equal(`"hello"`, string(encoder.Encode(nil, ptrOf("hello"))))
	should.Equal(`"\nhello中文"`, string(encoder.Encode(nil, ptrOf("\nhello中文"))))
	should.Equal(`"\nhello中文h\nello"`, string(encoder.Encode(nil, ptrOf("\nhello中文h\nello"))))
	should.Equal(`"\nhello中文h\nello\t"`, string(encoder.Encode(nil, ptrOf("\nhello中文h\nello\t"))))
}
