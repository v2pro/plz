package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"reflect"
	"io"
	"github.com/v2pro/plz/msgfmt/jsonfmt"
)

func Test_slice_of_empty_interface(t *testing.T) {
	should := require.New(t)
	encoder := jsonfmt.EncoderOf(reflect.TypeOf(([]interface{})(nil)))
	should.Equal("[1,null,3]", string(encoder.Encode(nil,nil, jsonfmt.PtrOf([]interface{}{
		1, nil, 3,
	}))))
}

type TestCloser int

func (closer TestCloser) Close() error {
	return nil
}

func Test_slice_of_non_empty_interface(t *testing.T) {
	should := require.New(t)
	encoder := jsonfmt.EncoderOf(reflect.TypeOf(([]io.Closer)(nil)))
	should.Equal("[1,null,3]", string(encoder.Encode(nil,nil, jsonfmt.PtrOf([]io.Closer{
		TestCloser(1), nil, TestCloser(3),
	}))))
}
