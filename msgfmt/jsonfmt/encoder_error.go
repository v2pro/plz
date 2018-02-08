package jsonfmt

import (
	"unsafe"
	"context"
)

type errorEncoder struct {
	sampleInterface emptyInterface
}

func (encoder *errorEncoder) Encode(ctx context.Context, space []byte, ptr unsafe.Pointer) []byte {
	errInterface := encoder.sampleInterface
	errInterface.word = ptr
	obj := *(*interface{})(unsafe.Pointer(&errInterface))
	space = append(space, '"')
	space = append(space, obj.(error).Error()...)
	space = append(space, '"')
	return space
}
