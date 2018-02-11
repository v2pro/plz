package reflect2

import (
	"unsafe"
	"reflect"
)


type unsafeType struct {
	reflect.Type
	rtype  unsafe.Pointer
	ptrRType unsafe.Pointer
}

func newUnsafeType(type1 reflect.Type) *unsafeType {
	return &unsafeType{
		Type: type1,
		rtype: toEface(type1).data,
		ptrRType: toEface(reflect.PtrTo(type1)).data,
	}
}

func (type2 *unsafeType) Type1() reflect.Type {
	return type2.Type
}

func (type2 *unsafeType) UnsafeNew() unsafe.Pointer {
	return unsafe_New(type2.rtype)
}

func (type2 *unsafeType) New() interface{} {
	return packEface(type2.ptrRType, type2.UnsafeNew())
}