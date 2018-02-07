package pickle

import (
	"unsafe"
	"reflect"
)

type rootEncoder struct {
	valType   reflect.Type
	signature uint64
	encoder   ValEncoder
}

func (encoder *rootEncoder) EncodeEmptyInterface(ptr unsafe.Pointer, stream *Stream) {
	stream.cursor = uintptr(len(stream.buf))
	valAsSlice := ptrAsBytes(int(encoder.valType.Size()), ptr)
	stream.buf = append(stream.buf, valAsSlice...)
	encoder.encoder.Encode(ptr, stream)
}

func (encoder *rootEncoder)  Encode(ptr unsafe.Pointer, stream *Stream) {
	encoder.encoder.Encode(ptr, stream)
}

func (encoder *rootEncoder) Signature() uint64 {
	return encoder.signature
}

func (encoder *rootEncoder) Type() reflect.Type {
	return encoder.valType
}

type rootDecoderWithCopy struct {
	valType   reflect.Type
	signature uint64
	decoder   ValDecoder
}

func (decoder *rootDecoderWithCopy) Signature() uint64 {
	return decoder.signature
}

func (decoder *rootDecoderWithCopy) Type() reflect.Type {
	return decoder.valType
}

func (decoder *rootDecoderWithCopy) DecodeEmptyInterface(ptr *emptyInterface, iter *Iterator) {
	iter.self = iter.allocator.Allocate(iter.objectSeq, iter.buf[16:16+decoder.Type().Size()])
	ptr.word = unsafe.Pointer(&iter.self[0])
	iter.cursor = iter.buf[16:]
	decoder.decoder.Decode(iter)
}

func (decoder *rootDecoderWithCopy) Decode(iter *Iterator) {
	decoder.decoder.Decode(iter)
}

type rootDecoderWithoutCopy struct {
	valType   reflect.Type
	signature uint64
	decoder   ValDecoder
}

func (decoder *rootDecoderWithoutCopy) Signature() uint64 {
	return decoder.signature
}

func (decoder *rootDecoderWithoutCopy) Type() reflect.Type {
	return decoder.valType
}

func (decoder *rootDecoderWithoutCopy) DecodeEmptyInterface(ptr *emptyInterface, iter *Iterator) {
	ptr.word = unsafe.Pointer(&iter.buf[16])
	iter.self = iter.buf[16:]
	iter.cursor = iter.buf[16:]
	decoder.decoder.Decode(iter)
}

func (decoder *rootDecoderWithoutCopy) Decode(iter *Iterator) {
	decoder.decoder.Decode(iter)
}
