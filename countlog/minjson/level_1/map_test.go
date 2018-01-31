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

func Test_map_of_ptr(t *testing.T) {
	should := require.New(t)
	one := 1
	encoder := minjson.EncoderOf(reflect.TypeOf(map[int]*int{1: &one}))
	should.Equal(`{"1":1}`, string(encoder.Encode(nil, minjson.PtrOf(map[int]*int{
		1: &one,
	}))))
}