package reflect2

import (
	"unsafe"
	"reflect"
)

type unsafeArrayType struct {
	unsafeType
	elemRType unsafe.Pointer
	elemSize uintptr
}

func newUnsafeArrayType(cfg *frozenConfig, type1 reflect.Type) *unsafeArrayType {
	return &unsafeArrayType{
		unsafeType: *newUnsafeType(cfg, type1),
		elemRType:  unpackEFace(type1.Elem()).data,
		elemSize:   type1.Elem().Size(),
	}
}

func (type2 *unsafeArrayType) Set(obj interface{}, index int, elem interface{}) {
	type2.UnsafeSet(unpackEFace(obj).data, index, unpackEFace(elem).data)
}

func (type2 *unsafeArrayType) UnsafeSet(obj unsafe.Pointer, index int, elem unsafe.Pointer) {
	elemPtr := arrayAt(obj, index, type2.elemSize, "i < s.Len")
	typedmemmove(type2.elemRType, elemPtr, elem)
}