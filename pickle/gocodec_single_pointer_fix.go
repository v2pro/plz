package pickle

import "unsafe"

type singlePointerEncoder struct {
	rootEncoder
}

func (encoder *singlePointerEncoder) EncodeEmptyInterface(ptr unsafe.Pointer, stream *Stream) {
	encoder.rootEncoder.EncodeEmptyInterface(unsafe.Pointer(&ptr), stream)
}

func (encoder *singlePointerEncoder) Encode(ptr unsafe.Pointer, stream *Stream) {
	encoder.rootEncoder.Encode(unsafe.Pointer(&ptr), stream)
}

type singlePointerDecoder struct {
	RootDecoder
}

func (decoder *singlePointerDecoder) DecodeEmptyInterface(ptr *emptyInterface, iter *Iterator) {
	decoder.RootDecoder.DecodeEmptyInterface(ptr, iter)
	ptr.word = *(*unsafe.Pointer)(ptr.word)
}