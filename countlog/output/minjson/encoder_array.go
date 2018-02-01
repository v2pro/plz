package minjson

import "unsafe"

type arrayEncoder struct {
	elemEncoder Encoder
	elemSize    uintptr
	length      int
}

func (encoder *arrayEncoder) Encode(space []byte, ptr unsafe.Pointer) []byte {
	space = append(space, '[')
	offset := uintptr(ptr)
	for i := 0; i < encoder.length; i++ {
		if i != 0 {
			space = append(space, ',')
		}
		space = encoder.elemEncoder.Encode(space, unsafe.Pointer(offset))
		offset += encoder.elemSize
	}
	space = append(space, ']')
	return space
}
