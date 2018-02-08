package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"reflect"
	"github.com/v2pro/plz/msgfmt/jsonfmt"
)

func Test_pointer(t *testing.T) {
	should := require.New(t)
	encoder := jsonfmt.EncoderOf(reflect.TypeOf((*int)(nil)))
	should.Equal("1", string(encoder.Encode(nil,nil, jsonfmt.PtrOf(1))))
	should.Equal("null", string(encoder.Encode(nil,nil, jsonfmt.PtrOf(nil))))
}