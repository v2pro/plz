package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"reflect"
	"github.com/v2pro/plz/msgfmt/njson"
)

func Test_slice(t *testing.T) {
	should := require.New(t)
	encoder := njson.EncoderOf(reflect.TypeOf([]int(nil)))
	should.Equal("[1,2,3]", string(encoder.Encode(nil, njson.PtrOf([]int{
		1, 2, 3,
	}))))
}