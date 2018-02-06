package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/v2pro/plz/countlog/pickle"
)

func Test_string(t *testing.T) {
	should := require.New(t)
	encoded, err := pickle.Marshal("hello")
	should.Nil(err)
	should.Equal([]byte{0x10, 0, 0, 0, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 'h', 'e', 'l', 'l', 'o'}, encoded[8:])
	decoded, err := pickle.ReadonlyConfig.Unmarshal(encoded, (*string)(nil))
	should.Nil(err)
	should.Equal("hello", *decoded.(*string))
	decoded, err = pickle.Unmarshal(encoded, (*string)(nil))
	should.Nil(err)
	should.Equal("hello", *decoded.(*string))
}