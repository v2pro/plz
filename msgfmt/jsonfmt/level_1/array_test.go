package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"reflect"
	"github.com/v2pro/plz/msgfmt/jsonfmt"
)

func Test_array(t *testing.T) {
	should := require.New(t)
	encoder := jsonfmt.EncoderOf(reflect.TypeOf([3]int{}))
	should.Equal("[1,2,3]", string(encoder.Encode(nil, nil, jsonfmt.PtrOf([3]int{
		1, 2, 3,
	}))))
}