package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"reflect"
	"io"
	"github.com/v2pro/plz/countlog/output/minjson"
)

func Test_map_of_number_key(t *testing.T) {
	should := require.New(t)
	encoder := minjson.EncoderOf(reflect.TypeOf(map[int]int{1: 1}))
	should.Equal(`{"1":1}`, string(encoder.Encode(nil, minjson.PtrOf(map[int]int{
		1: 1,
	}))))
}

func Test_map_of_string_key(t *testing.T) {
	should := require.New(t)
	encoder := minjson.EncoderOf(reflect.TypeOf(map[string]int{"hello": 1}))
	should.Equal(`{"hello":1}`, string(encoder.Encode(nil, minjson.PtrOf(map[string]int{
		"hello": 1,
	}))))
}

func Test_map_of_ptr_elem(t *testing.T) {
	should := require.New(t)
	one := 1
	encoder := minjson.EncoderOf(reflect.TypeOf(map[int]*int{1: &one}))
	should.Equal(`{"1":1}`, string(encoder.Encode(nil, minjson.PtrOf(map[int]*int{
		1: &one,
	}))))
}

func Test_map_of_interface_key(t *testing.T) {
	should := require.New(t)
	encoder := minjson.EncoderOf(reflect.TypeOf(map[interface{}]int{1: 1}))
	should.Equal(`{"1":1}`, string(encoder.Encode(nil, minjson.PtrOf(map[interface{}]int{
		1: 1,
	}))))
}

func Test_map_of_interface_elem(t *testing.T) {
	should := require.New(t)
	encoder := minjson.EncoderOf(reflect.TypeOf(map[int]interface{}{1: 1}))
	should.Equal(`{"1":1}`, string(encoder.Encode(nil, minjson.PtrOf(map[int]interface{}{
		1: 1,
	}))))
}

func Test_map_of_non_empty_interface_value(t *testing.T) {
	should := require.New(t)
	encoder := minjson.EncoderOf(reflect.TypeOf(map[int]io.Closer{1: nil}))
	should.Equal(`{"1":1}`, string(encoder.Encode(nil, minjson.PtrOf(map[int]io.Closer{
		1: TestCloser(1),
	}))))
}