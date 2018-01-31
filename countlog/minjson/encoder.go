package minjson

import (
	"unsafe"
	"reflect"
)

var bytesType = reflect.TypeOf([]byte(nil))

type Encoder interface {
	Encode(space []byte, ptr unsafe.Pointer) []byte
}

func EncoderOf(valType reflect.Type) Encoder {
	return encoderOf("", valType)
}

func encoderOf(prefix string, valType reflect.Type) Encoder {
	if bytesType == valType {
		return &bytesEncoder{}
	}
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
	case reflect.Ptr:
		elemEncoder := encoderOf(prefix + " [ptrElem]", valType.Elem())
		return &pointerEncoder{elemEncoder:elemEncoder}
	case reflect.Slice:
		elemEncoder := encoderOf(prefix + " [sliceElem]", valType.Elem())
		return &sliceEncoder{
			elemEncoder: elemEncoder,
			elemSize: valType.Elem().Size(),
		}
	case reflect.Array:
		elemEncoder := encoderOf(prefix + " [sliceElem]", valType.Elem())
		return &arrayEncoder{
			elemEncoder: elemEncoder,
			elemSize: valType.Elem().Size(),
			length: valType.Len(),
		}
	}
	return nil
}

func PtrOf(val interface{}) unsafe.Pointer {
	return (*emptyInterface)(unsafe.Pointer(&val)).word
}

// emptyInterface is the header for an interface{} value.
type emptyInterface struct {
	typ  unsafe.Pointer
	word unsafe.Pointer
}