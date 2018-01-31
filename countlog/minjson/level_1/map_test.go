package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/v2pro/plz/countlog/minjson"
	"reflect"
)

func Test_map(t *testing.T) {
	should := require.New(t)
	encoder := minjson.EncoderOf(reflect.TypeOf(map[int]int{1: 1}))
	should.Equal(`{"1":1}`, string(encoder.Encode(nil, minjson.PtrOf(map[int]int{
		1: 1,
	}))))
}