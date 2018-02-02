package njson

import "unsafe"

type structEncoder struct {
	fields []structEncoderField
}

type structEncoderField struct {
	offset  uintptr
	prefix  string
	encoder Encoder
}

func (encoder *structEncoder) Encode(space []byte, ptr unsafe.Pointer) []byte {
	space = append(space, '{')
	offset := uintptr(ptr)
	for _, field := range encoder.fields {
		space = append(space, field.prefix...)
		space = field.encoder.Encode(space, unsafe.Pointer(offset+field.offset))
	}
	space = append(space, '}')
	return space
}
