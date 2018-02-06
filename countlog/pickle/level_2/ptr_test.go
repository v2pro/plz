package test

import (
	"testing"
	"github.com/v2pro/plz/countlog/pickle"
	"github.com/stretchr/testify/require"
)

func Test_ptr_of_string(t *testing.T) {
	should := require.New(t)
	obj := "hello"
	encoded, err := pickle.Marshal(&obj)
	should.Nil(err)
	should.Equal([]byte{
		0x8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x10, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x5, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x68, 0x65, 0x6c, 0x6c, 0x6f,
	}, encoded[8:])
	decoded, err := pickle.ReadonlyConfig.Unmarshal(encoded, (**string)(nil))
	should.Nil(err)
	should.Equal("hello", **decoded.(**string))
	decoded, err = pickle.Unmarshal(encoded, (**string)(nil))
	should.Nil(err)
	should.Equal("hello", **decoded.(**string))
}

func Test_ptr_of_slice(t *testing.T) {
	should := require.New(t)
	obj := []byte("hello")
	encoded, err := pickle.Marshal(&obj)
	should.Nil(err)
	should.Equal([]byte{
		0x8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x18, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x5, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x5, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x68, 0x65, 0x6c, 0x6c, 0x6f,
	}, encoded[8:])
	decoded, err := pickle.ReadonlyConfig.Unmarshal(encoded, (**[]byte)(nil))
	should.Nil(err)
	should.Equal("hello", string(**decoded.(**[]byte)))
	decoded, err = pickle.Unmarshal(encoded, (**[]byte)(nil))
	should.Nil(err)
	should.Equal("hello", string(**decoded.(**[]byte)))
}
