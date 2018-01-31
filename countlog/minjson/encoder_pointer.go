package minjson

import (
	"unsafe"
)

type pointerEncoder struct {
	elemEncoder Encoder
}

func (encoder *pointerEncoder) Encode(space []byte, ptr unsafe.Pointer) []byte {
	if ptr == nil {
		return append(space, 'n', 'u', 'l', 'l')
	}
	return encoder.elemEncoder.Encode(space, *(*unsafe.Pointer)(ptr))
}
