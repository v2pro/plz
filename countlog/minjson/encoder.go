package minjson

import (
	"unsafe"
	"reflect"
)

type Encoder interface {
	Encode(space []byte, ptr unsafe.Pointer) []byte
}

func EncoderOf(valType reflect.Type) Encoder {
	return encoderOf("", valType)
}

func encoderOf(prefix string, valType reflect.Type) Encoder {
	switch valType.Kind() {
	case reflect.Int8:
		return &int8Encoder{}
	case reflect.Uint8:
		return &uint8Encoder{}
	case reflect.Int16:
		return &int16Encoder{}
	case reflect.Uint16:
		return &uint16Encoder{}
	case reflect.Int32:
		return &int32Encoder{}
	case reflect.Uint32:
		return &uint32Encoder{}
	case reflect.Int64, reflect.Int:
		return &int64Encoder{}
	case reflect.Uint64, reflect.Uint:
		return &uint64Encoder{}
	case reflect.Float64:
		return &lossyFloat64Encoder{}
	case reflect.Float32:
		return &lossyFloat32Encoder{}
	case reflect.String:
		return &stringEncoder{}
	}
	return nil
}