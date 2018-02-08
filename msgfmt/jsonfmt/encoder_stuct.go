package jsonfmt

import (
	"unsafe"
	"context"
)

type structEncoder struct {
	fields []structEncoderField
}

type structEncoderField struct {
	offset  uintptr
	prefix  string
	encoder Encoder
}

func (encoder *structEncoder) Encode(ctx context.Context, space []byte, ptr unsafe.Pointer) []byte {
	space = append(space, '{')
	offset := uintptr(ptr)
	for _, field := range encoder.fields {
		space = append(space, field.prefix...)
		space = field.encoder.Encode(ctx, space, unsafe.Pointer(offset+field.offset))
	}
	space = append(space, '}')
	return space
}
