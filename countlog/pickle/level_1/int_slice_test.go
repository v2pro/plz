package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/v2pro/plz/countlog/pickle"
)

func Test_int_slice(t *testing.T) {
	should := require.New(t)
	encoded, err := pickle.Marshal([]int{1, 2, 3})
	should.Nil(err)
	should.Equal([]byte{
		0x18, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0,
		1, 0, 0, 0, 0, 0, 0, 0,
		2, 0, 0, 0, 0, 0, 0, 0,
		3, 0, 0, 0, 0, 0, 0, 0}, encoded[8:])
	decoded, err := pickle.ReadonlyConfig.Unmarshal(encoded, (*[]int)(nil))
	should.Equal([]int{1, 2, 3}, *decoded.(*[]int))
	decoded, err = pickle.Unmarshal(encoded, (*[]int)(nil))
	should.Equal([]int{1, 2, 3}, *decoded.(*[]int))
}
