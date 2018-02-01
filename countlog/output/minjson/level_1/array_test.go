package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/v2pro/plz/countlog/minjson"
	"reflect"
)

func Test_array(t *testing.T) {
	should := require.New(t)
	encoder := minjson.EncoderOf(reflect.TypeOf([3]int{}))
	should.Equal("[1,2,3]", string(encoder.Encode(nil, minjson.PtrOf([3]int{
		1, 2, 3,
	}))))
}