package minjson

import (
	"unsafe"
	"reflect"
)

type onePtrInterfaceEncoder struct {
	valEncoder Encoder
}

func (encoder *onePtrInterfaceEncoder) Encode(space []byte, ptr unsafe.Pointer) []byte {
	return encoder.valEncoder.Encode(space, unsafe.Pointer(&ptr))
}

type emptyInterfaceEncoder struct {
}

func (encoder *emptyInterfaceEncoder) Encode(space []byte, ptr unsafe.Pointer) []byte {
	obj := *(*interface{})(ptr)
	return EncoderOf(reflect.TypeOf(obj)).Encode(space, PtrOf(obj))
}
