package pickle

import "unsafe"

type stringCodec struct {
	BaseCodec
}

func (codec *stringCodec) Encode(prStr unsafe.Pointer, stream *Stream) {
	pwStr := unsafe.Pointer(&stream.buf[stream.cursor])
	str := *(*string)(prStr)
	offset := uintptr(len(stream.buf)) - stream.cursor
	header := (*stringWritableHeader)(pwStr)
	header.Data = offset
	stream.buf = append(stream.buf, str...)
}

func (codec *stringCodec) Decode(iter *Iterator) {
	prStr := unsafe.Pointer(&iter.cursor[0])
	header := (*stringWritableHeader)(prStr)
	relOffset := header.Data
	pwStr := unsafe.Pointer(&iter.self[0])
	header = (*stringWritableHeader)(pwStr)
	header.Data = uintptr(unsafe.Pointer(&iter.cursor[relOffset]))
}

func (codec *stringCodec) HasPointer() bool {
	return true
}
