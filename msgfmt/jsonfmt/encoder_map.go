package jsonfmt

import (
	"unsafe"
	"reflect"
	"context"
)

type mapEncoder struct {
	keyEncoder      Encoder
	sampleInterface emptyInterface
}

func (encoder *mapEncoder) Encode(ctx context.Context, space []byte, ptr unsafe.Pointer) []byte {
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
		space = encoder.keyEncoder.Encode(ctx, space, unsafe.Pointer(&keyObj))
		elemObj := elem.Interface()
		space = EncoderOf(reflect.TypeOf(elemObj)).Encode(ctx, space, PtrOf(elemObj))
	}
	space = append(space, '}')
	return space
}

type mapNumberKeyEncoder struct {
	valEncoder Encoder
}

func (encoder *mapNumberKeyEncoder) Encode(ctx context.Context, space []byte, ptr unsafe.Pointer) []byte {
	space = append(space, '"')
	keyObj := *(*interface{})(ptr)
	space = encoder.valEncoder.Encode(ctx, space, PtrOf(keyObj))
	space = append(space, '"', ':')
	return space
}

type mapStringKeyEncoder struct {
	valEncoder Encoder
}

func (encoder *mapStringKeyEncoder) Encode(ctx context.Context, space []byte, ptr unsafe.Pointer) []byte {
	keyObj := *(*interface{})(ptr)
	space = encoder.valEncoder.Encode(ctx, space, PtrOf(keyObj))
	space = append(space, ':')
	return space
}

type mapInterfaceKeyEncoder struct {
	cfg *frozenConfig
	prefix string
}

func (encoder *mapInterfaceKeyEncoder) Encode(ctx context.Context, space []byte, ptr unsafe.Pointer) []byte {
	keyObj := *(*interface{})(ptr)
	return encoderOfMapKey(encoder.cfg, encoder.prefix, reflect.TypeOf(keyObj)).Encode(ctx, space, ptr)
}
