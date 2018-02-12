package jsonfmt

import (
	"unsafe"
	"reflect"
	"context"
	"github.com/v2pro/plz/reflect2"
)

type onePtrInterfaceEncoder struct {
	valEncoder Encoder
}

func (encoder *onePtrInterfaceEncoder) Encode(ctx context.Context, space []byte, ptr unsafe.Pointer) []byte {
	return encoder.valEncoder.Encode(ctx, space, unsafe.Pointer(&ptr))
}

type efaceEncoder struct {
}

func (encoder *efaceEncoder) Encode(ctx context.Context, space []byte, ptr unsafe.Pointer) []byte {
	obj := *(*interface{})(ptr)
	if obj == nil {
		return append(space, 'n', 'u', 'l', 'l')
	}
	return EncoderOf(reflect.TypeOf(obj)).Encode(ctx, space, PtrOf(obj))
}

type ifaceEncoder struct {
}

func (encoder *ifaceEncoder) Encode(ctx context.Context, space []byte, ptr unsafe.Pointer) []byte {
	obj := reflect2.UnsafeIFaceToEFace(ptr)
	if obj == nil {
		return append(space, 'n', 'u', 'l', 'l')
	}
	return EncoderOf(reflect.TypeOf(obj)).Encode(ctx, space, PtrOf(obj))
}
