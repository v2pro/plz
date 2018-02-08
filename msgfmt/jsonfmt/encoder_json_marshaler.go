package jsonfmt

import (
	"unsafe"
	"encoding/json"
	"context"
)

type jsonMarshalerEncoder struct {
	sampleInterface emptyInterface
}

func (encoder *jsonMarshalerEncoder) Encode(ctx context.Context, space []byte, ptr unsafe.Pointer) []byte {
	errInterface := encoder.sampleInterface
	errInterface.word = ptr
	obj := *(*interface{})(unsafe.Pointer(&errInterface))
	buf, err := obj.(json.Marshaler).MarshalJSON()
	if err != nil {
		space = append(space, '"')
		space = append(space, err.Error()...)
		space = append(space, '"')
		return space
	}
	space = append(space, buf...)
	return space
}