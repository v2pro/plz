package njson

import (
	"unsafe"
)

type errorEncoder struct {
	sampleInterface emptyInterface
}

func (encoder *errorEncoder) Encode(space []byte, pptr unsafe.Pointer) []byte {
	ptr := *(*unsafe.Pointer)(pptr)
	errInterface := encoder.sampleInterface
	errInterface.word = ptr
	obj := *(*interface{})(unsafe.Pointer(&errInterface))
	if obj == nil {
		return append(space, 'n', 'u', 'l', 'l')
	}
	space = append(space, '"')
	space = append(space, obj.(error).Error()...)
	space = append(space, '"')
	return space
}
