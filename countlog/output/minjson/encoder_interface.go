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
	if obj == nil {
		return append(space, 'n', 'u', 'l', 'l')
	}
	return EncoderOf(reflect.TypeOf(obj)).Encode(space, PtrOf(obj))
}

type nonEmptyInterfaceEncoder struct {
}

func (encoder *nonEmptyInterfaceEncoder) Encode(space []byte, ptr unsafe.Pointer) []byte {
	nonEmptyInterface := (*nonEmptyInterface)(ptr)
	var obj interface{}
	if nonEmptyInterface.itab != nil {
		e := (*emptyInterface)(unsafe.Pointer(&obj))
		e.typ = nonEmptyInterface.itab.typ
		e.word = nonEmptyInterface.word
	}
	if obj == nil {
		return append(space, 'n', 'u', 'l', 'l')
	}
	return EncoderOf(reflect.TypeOf(obj)).Encode(space, PtrOf(obj))
}

// emptyInterface is the header for an interface with method (not interface{})
type nonEmptyInterface struct {
	// see ../runtime/iface.go:/Itab
	itab *struct {
		ityp   unsafe.Pointer // static interface type
		typ    unsafe.Pointer // dynamic concrete type
		link   unsafe.Pointer
		bad    int32
		unused int32
		fun    [100000]unsafe.Pointer // method table
	}
	word unsafe.Pointer
}
