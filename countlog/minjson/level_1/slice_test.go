package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/v2pro/plz/countlog/minjson"
	"reflect"
)

func Test_slice(t *testing.T) {
	should := require.New(t)
	encoder := minjson.EncoderOf(reflect.TypeOf([]int(nil)))
	should.Equal("[1,2,3]", string(encoder.Encode(nil, minjson.PtrOf([]int{
		1, 2, 3,
	}))))
}