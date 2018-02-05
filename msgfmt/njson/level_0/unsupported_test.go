package test

import (
	"testing"
	"github.com/v2pro/plz/msgfmt/njson"
	"reflect"
	"github.com/stretchr/testify/require"
)

func Test_unsupported(t *testing.T) {
	should := require.New(t)
	encoder := njson.EncoderOf(reflect.TypeOf(make(chan int, 10)))
	output := encoder.Encode(nil, njson.PtrOf(make(chan int, 10)))
	should.Equal(`"can not encode chan int  to json"`, string(output))
}