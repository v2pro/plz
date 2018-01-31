package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/v2pro/plz/countlog/minjson"
	"reflect"
)

func Test_slice_of_empty_interface(t *testing.T) {
	should := require.New(t)
	encoder := minjson.EncoderOf(reflect.TypeOf(([]interface{})(nil)))
	should.Equal("[1,2,3]", string(encoder.Encode(nil, minjson.PtrOf([]interface{}{
		1, 2, 3,
	}))))

}
