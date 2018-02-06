package pickle

import "unsafe"

type singlePointerFix struct {
	rootEncoder
}

func (encoder *singlePointerFix) EncodeEmptyInterface(ptr unsafe.Pointer, stream *Stream) {
	encoder.rootEncoder.EncodeEmptyInterface(unsafe.Pointer(&ptr), stream)
}
