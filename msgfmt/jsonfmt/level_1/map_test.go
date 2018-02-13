package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"reflect"
	"io"
	"github.com/v2pro/plz/msgfmt/jsonfmt"
	"github.com/v2pro/plz/test"
	"github.com/v2pro/plz/countlog"
	"github.com/v2pro/plz/test/must"
	"encoding/json"
)

func Test_map_of_number_key(t *testing.T) {
	should := require.New(t)
	encoder := jsonfmt.EncoderOf(reflect.TypeOf(map[int]int{1: 1}))
	should.Equal(`{"1":1}`, string(encoder.Encode(nil,nil, jsonfmt.PtrOf(map[int]int{
		1: 1,
	}))))
}

func Test_map_of_string_key(t *testing.T) {
	should := require.New(t)
	encoder := jsonfmt.EncoderOf(reflect.TypeOf(map[string]int{"hello": 1}))
	should.Equal(`{"hello":1}`, string(encoder.Encode(nil,nil, jsonfmt.PtrOf(map[string]int{
		"hello": 1,
	}))))
}

func Test_map_of_ptr_elem(t *testing.T) {
	should := require.New(t)
	one := 1
	encoder := jsonfmt.EncoderOf(reflect.TypeOf(map[int]*int{1: &one}))
	should.Equal(`{"1":1}`, string(encoder.Encode(nil,nil, jsonfmt.PtrOf(map[int]*int{
		1: &one,
	}))))
}

func Test_map_of_interface_key(t *testing.T) {
	should := require.New(t)
	encoder := jsonfmt.EncoderOf(reflect.TypeOf(map[interface{}]int{1: 1}))
	should.Equal(`{"1":1}`, string(encoder.Encode(nil,nil, jsonfmt.PtrOf(map[interface{}]int{
		1: 1,
	}))))
}

func Test_map_of_interface_elem(t *testing.T) {
	should := require.New(t)
	encoder := jsonfmt.EncoderOf(reflect.TypeOf(map[int]interface{}{1: 1}))
	should.Equal(`{"1":1}`, string(encoder.Encode(nil,nil, jsonfmt.PtrOf(map[int]interface{}{
		1: 1,
	}))))
}

func Test_map_of_non_empty_interface_value(t *testing.T) {
	should := require.New(t)
	encoder := jsonfmt.EncoderOf(reflect.TypeOf(map[int]io.Closer{1: nil}))
	should.Equal(`{"1":1}`, string(encoder.Encode(nil,nil, jsonfmt.PtrOf(map[int]io.Closer{
		1: TestCloser(1),
	}))))
}

func Test_map(t *testing.T) {
	t.Run("map string to eface", test.Case(func(ctx *countlog.Context) {
		must.JsonEqual(`{
			"hello": 1,
			"world": "yes"
		}`, jsonfmt.MarshalToString(map[string]interface{}{
			"hello": 1,
			"world": "yes",
		}))
	}))
}

func Benchmark_map_unsafe(b *testing.B) {
	encoder := jsonfmt.EncoderOf(reflect.TypeOf(map[string]int{}))
	m := map[string]int {
		"hello": 1,
		"world": 3,
	}
	b.ReportAllocs()
	b.ResetTimer()
	space := []byte(nil)
	for i := 0; i < b.N; i++ {
		space = encoder.Encode(nil, space[:0], jsonfmt.PtrOf(m))
	}
}

func Benchmark_map_safe(b *testing.B) {
	m := map[string]int {
		"hello": 1,
		"world": 3,
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		json.Marshal(m)
	}
}