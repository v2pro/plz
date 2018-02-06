package pickle

import (
	"fmt"
	"reflect"
	"unsafe"
)

type Stream struct {
	cfg *frozenConfig
	// there are two pointers
	// buf + cursor => the input of encoder
	// buf + len(buf) => the output of encoder
	buf    []byte
	cursor uintptr
	Error  error
}

func (cfg *frozenConfig) NewStream(buf []byte) *Stream {
	return &Stream{cfg: cfg, buf: buf}
}

func (stream *Stream) Reset(buf []byte) {
	stream.buf = buf
	stream.cursor = 0
}

func (stream *Stream) Marshal(val interface{}) uint32 {
	valType := reflect.TypeOf(val)
	encoder, err := encoderOfType(stream.cfg, valType)
	if err != nil {
		stream.ReportError("EncodeVal", err)
		return 0
	}
	baseCursor := len(stream.buf)
	stream.buf = append(stream.buf, []byte{
		0, 0, 0, 0, // size
		0, 0, 0, 0, 0, 0, 0, 0, // signature
	}...)
	encoder.EncodeEmptyInterface(ptrOfEmptyInterface(val), stream)
	if stream.Error != nil {
		return 0
	}
	pSize := unsafe.Pointer(&stream.buf[baseCursor])
	size := uint32(len(stream.buf) - baseCursor)
	*(*uint32)(pSize) = size
	pSig := unsafe.Pointer(&stream.buf[baseCursor+4])
	if stream.cfg.useSignature {
		*(*uint64)(pSig) = encoder.Signature()
	} else {
		*(*uint64)(pSig) = uint64(uintptr(ptrOfEmptyInterface(valType)))
	}
	return size
}

func (stream *Stream) Buffer() []byte {
	return stream.buf
}

func (stream *Stream) ReportError(operation string, err error) {
	if stream.Error != nil {
		return
	}
	stream.Error = fmt.Errorf("%s: %s", operation, err)
}
