package njson

import (
	"unsafe"
)

type pointerEncoder struct {
	elemEncoder Encoder
}

func (encoder *pointerEncoder) Encode(space []byte, ptr unsafe.Pointer) []byte {
	ptrTo := *(*unsafe.Pointer)(ptr)
	if ptrTo == nil {
		return append(space, 'n', 'u', 'l', 'l')
	}
	return encoder.elemEncoder.Encode(space, ptrTo)
}
