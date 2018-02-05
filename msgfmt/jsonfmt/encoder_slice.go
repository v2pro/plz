package jsonfmt

import "unsafe"

type sliceEncoder struct {
	elemEncoder Encoder
	elemSize    uintptr
}

func (encoder *sliceEncoder) Encode(space []byte, ptr unsafe.Pointer) []byte {
	slice := (*sliceHeader)(ptr)
	space = append(space, '[')
	offset := uintptr(slice.Data)
	for i := 0; i < slice.Len; i++ {
		if i != 0 {
			space = append(space, ',')
		}
		space = encoder.elemEncoder.Encode(space, unsafe.Pointer(offset))
		offset += encoder.elemSize
	}
	space = append(space, ']')
	return space
}
