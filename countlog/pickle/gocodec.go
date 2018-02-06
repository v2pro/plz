package pickle

import (
	"unsafe"
	"reflect"
	"sync"
)

type ObjectSeq uint64

type Allocator interface {
	Allocate(ObjectSeq, []byte) []byte
}

type DefaultAllocator struct {
}

func (allocator *DefaultAllocator) Allocate(objectSeq ObjectSeq, original []byte) []byte {
	return append([]byte(nil), original...)
}

var defaultAllocator = &DefaultAllocator{}

type Config struct {
	UseSignature   bool
	ReadonlyDecode bool
}

type API interface {
	Marshal(val interface{}) ([]byte, error)
	Unmarshal(buf []byte) (interface{}, error)
	UnmarshalCandidates(buf []byte, candidatePointers ...interface{}) (interface{}, error)
	NewIterator(buf []byte) *Iterator
	NewStream(buf []byte) *Stream
}

type ValEncoder interface {
	Encode(ptr unsafe.Pointer, stream *Stream)
	Type() reflect.Type
	IsNoop() bool
	Signature() uint64
}

type RootEncoder interface {
	Type() reflect.Type
	Signature() uint64
	EncodeEmptyInterface(ptr unsafe.Pointer, stream *Stream)
}

type ValDecoder interface {
	Decode(iter *Iterator)
	Type() reflect.Type
	IsNoop() bool
	Signature() uint64
	HasPointer() bool
}

type RootDecoder interface {
	Type() reflect.Type
	Signature() uint64
	DecodeEmptyInterface(ptr *emptyInterface, iter *Iterator)
}

type frozenConfig struct {
	useSignature   bool
	readonlyDecode bool
	allocator      Allocator
	decoderCache   *sync.Map
	encoderCache   *sync.Map
}

func (cfg Config) Froze() API {
	api := &frozenConfig{
		useSignature:   cfg.UseSignature,
		readonlyDecode: cfg.ReadonlyDecode,
		decoderCache:   &sync.Map{},
		encoderCache:   &sync.Map{},
	}
	return api
}

var ReadonlyConfig = Config{ReadonlyDecode:true}.Froze()
var DefaultConfig = Config{}.Froze()

func Marshal(obj interface{}) ([]byte, error) {
	return DefaultConfig.Marshal(obj)
}

func Unmarshal(buf []byte) (interface{}, error) {
	return DefaultConfig.Unmarshal(buf)
}

func UnmarshalCandidates(buf []byte, candidatePointers ...interface{}) (interface{}, error) {
	return DefaultConfig.UnmarshalCandidates(buf, candidatePointers...)
}

func NewIterator(buf []byte) *Iterator {
	return DefaultConfig.NewIterator(buf)
}

func NewStream(buf []byte) *Stream {
	return DefaultConfig.NewStream(buf)
}

func (cfg *frozenConfig) Marshal(val interface{}) ([]byte, error) {
	stream := cfg.NewStream(nil)
	stream.Marshal(val)
	return stream.Buffer(), stream.Error
}

func (cfg *frozenConfig) Unmarshal(buf []byte) (interface{}, error) {
	iter := cfg.NewIterator(buf)
	val := iter.Unmarshal()
	return val, iter.Error
}

func (cfg *frozenConfig) UnmarshalCandidates(buf []byte, candidatePointers ...interface{}) (interface{}, error) {
	iter := cfg.NewIterator(buf)
	val := iter.UnmarshalCandidates(candidatePointers...)
	return val, iter.Error
}
