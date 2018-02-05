package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"reflect"
	"github.com/v2pro/plz/msgfmt/jsonfmt"
	"errors"
)

func Test_error(t *testing.T) {
	should := require.New(t)
	encoder := jsonfmt.EncoderOf(reflect.TypeOf(errors.New("hello")))
	should.Equal(`"hello"`, string(encoder.Encode(nil, jsonfmt.PtrOf(errors.New("hello")))))
}