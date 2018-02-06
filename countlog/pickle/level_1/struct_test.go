package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/v2pro/plz/countlog/pickle"
)

func Test_single_ptr_in_struct(t *testing.T) {
	should := require.New(t)
	type TestObject struct {
		Field1 *uint8
	}
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

func Test_two_ptrs_in_struct(t *testing.T) {
	should := require.New(t)
	type TestObject struct {
		Field1 *uint8
		Field2 *uint8
	}
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

func Test_slice_in_struct(t *testing.T) {
	should := require.New(t)
	type TestObject struct {
		Field1 uint
		Field2 []uint
	}
	obj := TestObject{Field1: 1, Field2: []uint{2,3}}
	encoded, err := pickle.Marshal(obj)
	should.Nil(err)
	should.Equal([]byte{
		0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x18, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	}, encoded[8:])
	decoded, err := pickle.ReadonlyConfig.Unmarshal(encoded, (*TestObject)(nil))
	should.Nil(err)
	should.Equal(obj, *decoded.(*TestObject))
	decoded, err = pickle.Unmarshal(encoded, (*TestObject)(nil))
	should.Nil(err)
	should.Equal(obj, *decoded.(*TestObject))
}

func Test_multiple_struct(t *testing.T) {
	should := require.New(t)
	type TestObj struct {
		Field []int
	}
	stream := pickle.NewStream(nil)
	stream.Marshal(TestObj{[]int{1}})
	stream.Marshal(TestObj{[]int{2}})
	should.Nil(stream.Error)
	iter := pickle.NewIterator(stream.Buffer())
	obj := iter.Unmarshal((*TestObj)(nil))
	should.Nil(iter.Error)
	should.Equal([]int{1}, obj.(*TestObj).Field)
	obj = iter.Unmarshal((*TestObj)(nil))
	should.Nil(iter.Error)
	should.Equal([]int{2}, obj.(*TestObj).Field)
}
