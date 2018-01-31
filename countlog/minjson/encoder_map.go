package minjson

import (
	"unsafe"
	"reflect"
)

type mapEncoder struct {
	keyEncoder      Encoder
	sampleInterface emptyInterface
}

func (encoder *mapEncoder) Encode(space []byte, ptr unsafe.Pointer) []byte {
	mapInterface := encoder.sampleInterface
	mapInterface.word = ptr
	mapObj := *(*interface{})(unsafe.Pointer(&mapInterface))
	mapVal := reflect.ValueOf(mapObj)
	if mapVal.IsNil() {
		return append(space, 'n', 'u', 'l', 'l')
	}
	keys := mapVal.MapKeys()
	space = append(space, '{')
	for i, key := range keys {
		if i != 0 {
			space = append(space, ',')
		}
		elem := mapVal.MapIndex(key)
		keyObj := key.Interface()
		space = encoder.keyEncoder.Encode(space, unsafe.Pointer(&keyObj))
		elemObj := elem.Interface()
		space = EncoderOf(reflect.TypeOf(elemObj)).Encode(space, PtrOf(elemObj))
	}
	space = append(space, '}')
	return space
}

type mapNumberKeyEncoder struct {
	valEncoder Encoder
}

func (encoder *mapNumberKeyEncoder) Encode(space []byte, ptr unsafe.Pointer) []byte {
	space = append(space, '"')
	keyObj := *(*interface{})(ptr)
	space = encoder.valEncoder.Encode(space, PtrOf(keyObj))
	space = append(space, '"', ':')
	return space
}

type mapStringKeyEncoder struct {
	valEncoder Encoder
}

func (encoder *mapStringKeyEncoder) Encode(space []byte, ptr unsafe.Pointer) []byte {
	keyObj := *(*interface{})(ptr)
	space = encoder.valEncoder.Encode(space, PtrOf(keyObj))
	space = append(space, ':')
	return space
}

type mapInterfaceKeyEncoder struct {
	prefix string
}

func (encoder *mapInterfaceKeyEncoder) Encode(space []byte, ptr unsafe.Pointer) []byte {
	keyObj := *(*interface{})(ptr)
	return encoderOfMapKey("", reflect.TypeOf(keyObj)).Encode(space, ptr)
}
