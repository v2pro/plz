package functional

import (
	"testing"
	"unsafe"
	"github.com/stretchr/testify/require"
)

func Test_iterate_elements(t *testing.T) {
	should := require.New(t)
	val := []int{1, 2, 3}
	fiz := getFp(val)
	counter := 0
	fiz.iterateElements(toPointer(val), func(elem unsafe.Pointer) bool {
		counter++
		return true
	})
	should.Equal(3, counter)
}
