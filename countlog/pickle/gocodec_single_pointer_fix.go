package pickle

import "unsafe"

type singlePointerEncoder struct {
	rootEncoder
}

func (encoder *singlePointerEncoder) EncodeEmptyInterface(ptr unsafe.Pointer, stream *Stream) {
	encoder.rootEncoder.EncodeEmptyInterface(unsafe.Pointer(&ptr), stream)
}

type singlePointerDecoder struct {
	RootDecoder
}

func (encoder *singlePointerDecoder) DecodeEmptyInterface(ptr *emptyInterface, iter *Iterator) {
	encoder.RootDecoder.DecodeEmptyInterface(ptr, iter)
	ptr.word = *(*unsafe.Pointer)(ptr.word)
}
