package pickle

import (
	"reflect"
	"fmt"
	"unsafe"
	"errors"
	"io"
)

type Iterator struct {
	allocator Allocator
	objectSeq ObjectSeq
	cfg       *frozenConfig
	buf       []byte
	self      []byte
	cursor    []byte
	Error     error
}

func (cfg *frozenConfig) NewIterator(buf []byte) *Iterator {
	return &Iterator{cfg: cfg, buf: buf, allocator: defaultAllocator}
}

func (iter *Iterator) Reset(buf []byte) {
	iter.buf = buf
	iter.cursor = nil
	iter.Error = nil
}

func (iter *Iterator) NextSize() uint32 {
	if len(iter.buf) < 8 {
		return 0
	}
	return *(*uint32)(unsafe.Pointer(&iter.buf[0]))
}

func (iter *Iterator) Skip() []byte {
	size := iter.NextSize()
	skipped := iter.buf[:size]
	iter.buf = iter.buf[size:]
	return skipped
}

func (iter *Iterator) CopyThenUnmarshal() interface{} {
	size := iter.NextSize()
	if size == 0 {
		iter.Error = io.EOF
		return nil
	}
	copied := iter.allocator.Allocate(iter.objectSeq, iter.buf[:size])
	nextBuf := iter.buf[size:]
	iter.Reset(copied)
	result := iter.Unmarshal()
	iter.Reset(nextBuf)
	return result
}

func (iter *Iterator) CopyThenUnmarshalCandidates(candidatePointers ...interface{}) interface{} {
	size := iter.NextSize()
	if size == 0 {
		iter.Error = io.EOF
		return nil
	}
	copied := iter.allocator.Allocate(iter.objectSeq, iter.buf[:size])
	nextBuf := iter.buf[size:]
	iter.Reset(copied)
	result := iter.UnmarshalCandidates(candidatePointers...)
	iter.Reset(nextBuf)
	return result
}

func (iter *Iterator) ObjectSeq(objectSeq ObjectSeq) {
	iter.objectSeq = objectSeq
}

func (iter *Iterator) Allocator(allocator Allocator) {
	iter.allocator = allocator
}

var sampleInterface emptyInterface

func init() {
	var typ interface{} = reflect.TypeOf("")
	sampleInterface = *(*emptyInterface)(unsafe.Pointer(&typ))
}

func (iter *Iterator) Unmarshal() interface{} {
	size := iter.NextSize()
	if size == 0 {
		iter.Error = io.EOF
		return nil
	}
	defer func() {
		recovered := recover()
		if recovered != nil {
			iter.ReportError("Unmarshal", fmt.Errorf("%v", recovered))
		}
	}()
	nextBuf := iter.buf[size:]
	typeId := *(*uint64)(unsafe.Pointer(&iter.buf[4]))
	valTypeInterface := sampleInterface
	valTypeInterface.word = unsafe.Pointer(uintptr(typeId))
	valType := (*(*interface{})(unsafe.Pointer(&valTypeInterface))).(reflect.Type)
	if valType == nil {
		iter.ReportError("Unmarshal", fmt.Errorf("unknown type: %v", typeId))
		return nil
	}
	var decoder RootDecoder
	val := reflect.New(valType).Elem().Interface()
	decoder, err := decoderOfType(iter.cfg, valType)
	if err != nil {
		iter.ReportError("DecodeVal", err)
		return nil
	}
	decoder.DecodeEmptyInterface((*emptyInterface)(unsafe.Pointer(&val)), iter)
	iter.buf = nextBuf
	return val
}

func (iter *Iterator) UnmarshalCandidates(candidatePointers ...interface{}) interface{} {
	size := iter.NextSize()
	if size == 0 {
		iter.Error = io.EOF
		return nil
	}
	defer func() {
		recovered := recover()
		if recovered != nil {
			iter.ReportError("Unmarshal", fmt.Errorf("%v", recovered))
		}
	}()
	nextBuf := iter.buf[size:]
	sig := *(*uint64)(unsafe.Pointer(&iter.buf[4]))
	var decoder RootDecoder
	var val interface{}
	for _, candidatePointer := range candidatePointers {
		valType := reflect.TypeOf(candidatePointer).Elem()
		tryDecoder, err := decoderOfType(iter.cfg, valType)
		if err != nil {
			iter.ReportError("DecodeVal", err)
			return nil
		}
		if tryDecoder.Signature() == sig {
			decoder = tryDecoder
			val = candidatePointer
			break
		}
	}
	if decoder == nil {
		iter.ReportError("DecodeVal", errors.New("no decoder matches the signature"))
		return nil
	}
	decoder.DecodeEmptyInterface((*emptyInterface)(unsafe.Pointer(&val)), iter)
	iter.buf = nextBuf
	return val
}

func (iter *Iterator) ReportError(operation string, err error) {
	if iter.Error != nil {
		return
	}
	iter.Error = fmt.Errorf("%s: %s", operation, err)
}

func (iter *Iterator) Buffer() []byte {
	return iter.buf
}
