package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"reflect"
	"github.com/v2pro/plz/nfmt/njson"
	"errors"
)

func Test_error(t *testing.T) {
	should := require.New(t)
	encoder := njson.EncoderOf(reflect.TypeOf(errors.New("hello")))
	should.Equal(`"hello"`, string(encoder.Encode(nil, njson.PtrOf(errors.New("hello")))))
}