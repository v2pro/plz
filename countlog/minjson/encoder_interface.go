package minjson

import "unsafe"

type onePtrInterfaceEncoder struct {
	valEncoder Encoder
}

func (encoder *onePtrInterfaceEncoder) Encode(space []byte, ptr unsafe.Pointer) []byte {
	return encoder.valEncoder.Encode(space, unsafe.Pointer(&ptr))
}