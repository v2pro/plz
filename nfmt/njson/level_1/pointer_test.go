package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"reflect"
	"github.com/v2pro/plz/nfmt/njson"
)

func Test_pointer(t *testing.T) {
	should := require.New(t)
	encoder := njson.EncoderOf(reflect.TypeOf((*int)(nil)))
	should.Equal("1", string(encoder.Encode(nil, njson.PtrOf(1))))
	should.Equal("null", string(encoder.Encode(nil, njson.PtrOf(nil))))
}