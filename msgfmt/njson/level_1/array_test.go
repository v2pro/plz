package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"reflect"
	"github.com/v2pro/plz/msgfmt/njson"
)

func Test_array(t *testing.T) {
	should := require.New(t)
	encoder := njson.EncoderOf(reflect.TypeOf([3]int{}))
	should.Equal("[1,2,3]", string(encoder.Encode(nil, njson.PtrOf([3]int{
		1, 2, 3,
	}))))
}