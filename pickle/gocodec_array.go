package pickle

import "unsafe"

type arrayEncoder struct {
	BaseCodec
	arrayLength int
	elementSize uintptr
	elemEncoder ValEncoder
}

func (encoder *arrayEncoder) Encode(prArray unsafe.Pointer, stream *Stream) {
	if encoder.IsNoop() {
		return
	}
	cursor := stream.cursor
	prElem := uintptr(prArray)
	for i := 0; i < encoder.arrayLength; i++ {
		stream.cursor = cursor // stream.cursor will change in the elemEncoder
		encoder.elemEncoder.Encode(unsafe.Pointer(prElem), stream)
		cursor = cursor + encoder.elementSize
		prElem = prElem + encoder.elementSize
	}
}

func (encoder *arrayEncoder) IsNoop() bool {
	return encoder.elemEncoder == nil
}

type arrayDecoderWithoutPointer struct {
	BaseCodec
	arrayLength int
	elementSize uintptr
	elemDecoder ValDecoder
}

func (decoder *arrayDecoderWithoutPointer) Decode(iter *Iterator) {
	if decoder.IsNoop() {
		return
	}
	cursor := iter.cursor
	for i := 0; i < decoder.arrayLength; i++ {
		iter.cursor = cursor // iter.cursor will change in elemDecoder
		decoder.elemDecoder.Decode(iter)
		cursor = cursor[decoder.elementSize:]
	}
}

func (decoder *arrayDecoderWithoutPointer) IsNoop() bool {
	return decoder.elemDecoder == nil
}

type arrayDecoderWithPointer struct {
	BaseCodec
	arrayLength int
	elementSize uintptr
	elemDecoder ValDecoder
}

func (decoder *arrayDecoderWithPointer) Decode(iter *Iterator) {
	if decoder.IsNoop() {
		return
	}
	cursor := iter.cursor
	self := iter.self
	for i := 0; i < decoder.arrayLength; i++ {
		iter.cursor = cursor // iter.cursor will change in elemDecoder
		iter.self = self
		decoder.elemDecoder.Decode(iter)
		cursor = cursor[decoder.elementSize:]
		self = self[decoder.elementSize:]
	}
}

func (decoder *arrayDecoderWithPointer) IsNoop() bool {
	return decoder.elemDecoder == nil
}

func (decoder *arrayDecoderWithPointer) HasPointer() bool {
	return true
}
