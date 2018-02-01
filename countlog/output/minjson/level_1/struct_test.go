package test

import (
	"testing"
	"reflect"
	"github.com/stretchr/testify/require"
	"github.com/v2pro/plz/countlog/output/minjson"
)

func Test_struct(t *testing.T) {
	should := require.New(t)
	type TestObject struct {
		Field1 string
		Field2 int `json:"field_2"`
	}
	encoder := minjson.EncoderOf(reflect.TypeOf(TestObject{}))
	output := encoder.Encode(nil, minjson.PtrOf(TestObject{"hello", 100}))
	should.Equal(`{"Field1":"hello","field_2":100}`, string(output))
}
