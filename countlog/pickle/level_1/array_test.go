package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/v2pro/plz/countlog/pickle"
)

func Test_single_ptr_int_array(t *testing.T) {
	should := require.New(t)
	type TestObject [1]*uint8
	one := uint8(1)
	obj := TestObject{&one}
	encoded, err := pickle.Marshal(obj)
	should.Nil(err)
	should.Equal([]byte{
		0x8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x1,
	}, encoded[8:])
	decoded, err := pickle.ReadonlyConfig.Unmarshal(encoded, (*TestObject)(nil))
	should.Nil(err)
	should.Equal(obj, *decoded.(*TestObject))
	decoded, err = pickle.Unmarshal(encoded, (*TestObject)(nil))
	should.Nil(err)
	should.Equal(obj, *decoded.(*TestObject))
}

func Test_two_ptrs_in_array(t *testing.T) {
	should := require.New(t)
	type TestObject [2]*uint8
	one := uint8(1)
	obj := TestObject{&one, &one}
	encoded, err := pickle.Marshal(obj)
	should.Nil(err)
	should.Equal([]byte{
		0x10, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x9, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		1,
		1,
	}, encoded[8:])
	decoded, err := pickle.ReadonlyConfig.Unmarshal(encoded, (*TestObject)(nil))
	should.Nil(err)
	should.Equal(obj, *decoded.(*TestObject))
	decoded, err = pickle.Unmarshal(encoded, (*TestObject)(nil))
	should.Nil(err)
	should.Equal(obj, *decoded.(*TestObject))
}