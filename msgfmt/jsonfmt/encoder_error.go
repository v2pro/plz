package jsonfmt

import (
	"unsafe"
)

type errorEncoder struct {
	sampleInterface emptyInterface
}

func (encoder *errorEncoder) Encode(space []byte, ptr unsafe.Pointer) []byte {
	errInterface := encoder.sampleInterface
	errInterface.word = ptr
	obj := *(*interface{})(unsafe.Pointer(&errInterface))
	space = append(space, '"')
	space = append(space, obj.(error).Error()...)
	space = append(space, '"')
	return space
}
