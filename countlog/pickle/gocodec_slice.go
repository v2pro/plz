package pickle

import "unsafe"

type sliceEncoder struct {
	BaseCodec
	elemSize    int
	elemEncoder ValEncoder
}

func (encoder *sliceEncoder) Encode(prSlice unsafe.Pointer, stream *Stream) {
	rHeader := (*sliceReadonlyHeader)(prSlice)
	if rHeader.Len == 0 {
		return
	}
	pwSlice := unsafe.Pointer(&stream.buf[stream.cursor])
	wHeader := (*sliceWritableHeader)(pwSlice)
	wHeader.Cap = rHeader.Len
	byteSlice := ptrAsBytes(encoder.elemSize*rHeader.Len, rHeader.Data)
	// replace actual pointer with relative offset
	wHeader.Data = uintptr(len(stream.buf)) - stream.cursor
	stream.cursor = uintptr(len(stream.buf)) // start of the bytes
	stream.buf = append(stream.buf, byteSlice...)
	if encoder.elemEncoder != nil {
		endCursor := uintptr(len(stream.buf)) // end of the bytes
		cursor := stream.cursor
		prElem := uintptr(rHeader.Data)
		for ; cursor < endCursor; cursor += uintptr(encoder.elemSize) {
			stream.cursor = cursor
			encoder.elemEncoder.Encode(unsafe.Pointer(prElem), stream)
			prElem += uintptr(encoder.elemSize)
		}
	}
}

type sliceDecoderWithoutCopy struct {
	BaseCodec
	elemSize    int
	elemDecoder ValDecoder
}

func (decoder *sliceDecoderWithoutCopy) Decode(iter *Iterator) {
	pwSlice := unsafe.Pointer(&iter.self[0])
	header := (*sliceWritableHeader)(pwSlice)
	if header.Len == 0 {
		return
	}
	relOffset := header.Data
	header.Data = uintptr(unsafe.Pointer(&iter.cursor[relOffset]))
	if decoder.elemDecoder != nil {
		cursor := iter.cursor[relOffset:]
		for i := 0; i < header.Len; i++ {
			if i > 0 {
				cursor = cursor[decoder.elemSize:]
			}
			iter.cursor = cursor
			iter.self = iter.cursor
			decoder.elemDecoder.Decode(iter)
		}
	}
}

func (decoder *sliceDecoderWithoutCopy) HasPointer() bool {
	return true
}

type sliceDecoderWithCopy struct {
	BaseCodec
	elemSize    int
	elemDecoder ValDecoder
}

func (decoder *sliceDecoderWithCopy) Decode(iter *Iterator) {
	pwSlice := unsafe.Pointer(&iter.self[0])
	header := (*sliceWritableHeader)(pwSlice)
	if header.Len == 0 {
		return
	}
	relOffset := header.Data
	cursor := iter.cursor[relOffset:]
	copied := iter.allocator.Allocate(iter.objectSeq, cursor[:decoder.elemSize*header.Len])
	header.Data = uintptr(unsafe.Pointer(&copied[0]))
	for i := 0; i < header.Len; i++ {
		if i > 0 {
			cursor = cursor[decoder.elemSize:]
			copied = copied[decoder.elemSize:]
		}
		iter.cursor = cursor
		iter.self = copied
		decoder.elemDecoder.Decode(iter)
	}
}

func (decoder *sliceDecoderWithCopy) HasPointer() bool {
	return true
}
