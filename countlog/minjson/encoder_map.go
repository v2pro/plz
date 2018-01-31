package minjson

import (
	"unsafe"
	"reflect"
)

type mapEncoder struct {
	keyEncoder      Encoder
	elemEncoder     Encoder
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
		space = encoder.keyEncoder.Encode(space, PtrOf(key.Interface()))
		space = append(space, ':')
		space = encoder.elemEncoder.Encode(space, PtrOf(elem.Interface()))
	}
	space = append(space, '}')
	return space
}

type mapKeyEncoder struct {
	valEncoder Encoder
}

func (encoder *mapKeyEncoder) Encode(space []byte, ptr unsafe.Pointer) []byte {
	space = append(space, '"')
	space = encoder.valEncoder.Encode(space, ptr)
	space = append(space, '"')
	return space
}
