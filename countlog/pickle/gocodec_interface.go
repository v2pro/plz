package pickle

import "unsafe"

type interfaceEncoder struct {
	BaseCodec
}

func (encoder *interfaceEncoder) Encode(prPointer unsafe.Pointer, stream *Stream) {
}