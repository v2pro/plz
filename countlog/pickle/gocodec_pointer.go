package pickle

import (
	"unsafe"
)

type pointerEncoder struct {
	BaseCodec
	elemEncoder ValEncoder
}

func (encoder *pointerEncoder) Encode(prPointer unsafe.Pointer, stream *Stream) {
	ptr := *(*unsafe.Pointer)(prPointer)
	if uintptr(ptr) == 0 {
		return
	}
	valAsBytes := ptrAsBytes(int(encoder.elemEncoder.Type().Size()), ptr)
	pwPointer := unsafe.Pointer(&stream.buf[stream.cursor])
	*(*uintptr)(pwPointer) = uintptr(len(stream.buf)) - stream.cursor
	stream.cursor = uintptr(len(stream.buf))
	stream.buf = append(stream.buf, valAsBytes...)
	encoder.elemEncoder.Encode(ptr, stream)
}

type pointerDecoderWithoutCopy struct {
	BaseCodec
	elemDecoder ValDecoder
}

func (decoder *pointerDecoderWithoutCopy) Decode(iter *Iterator) {
	pPtr := unsafe.Pointer(&iter.cursor[0])
	relOffset := *(*uintptr)(pPtr)
	if relOffset == 0 {
		return
	}
	iter.cursor = iter.cursor[relOffset:]
	*(*uintptr)(unsafe.Pointer(&iter.self[0])) = uintptr(unsafe.Pointer(&iter.cursor[0]))
	iter.self = iter.cursor
	decoder.elemDecoder.Decode(iter)
}

func (decoder *pointerDecoderWithoutCopy) HasPointer() bool {
	return true
}

type pointerDecoderWithCopy struct {
	BaseCodec
	elemDecoder ValDecoder
}

func (decoder *pointerDecoderWithCopy) Decode(iter *Iterator) {
	pPtr := unsafe.Pointer(&iter.cursor[0])
	relOffset := *(*uintptr)(pPtr)
	if relOffset == 0 {
		return
	}
	iter.cursor = iter.cursor[relOffset:]
	copied := iter.allocator.Allocate(iter.objectSeq, iter.cursor[:decoder.elemDecoder.Type().Size()])
	*(*uintptr)(unsafe.Pointer(&iter.self[0])) = uintptr(unsafe.Pointer(&copied[0]))
	iter.self = copied
	decoder.elemDecoder.Decode(iter)
}

func (decoder *pointerDecoderWithCopy) HasPointer() bool {
	return true
}
