package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/v2pro/plz/msgfmt/jsonfmt"
	"reflect"
)
func Test_bool(t *testing.T) {
	should := require.New(t)
	encoder := jsonfmt.EncoderOf(reflect.TypeOf(true))
	should.Equal("true", string(encoder.Encode(nil, nil, jsonfmt.PtrOf(true))))
}
