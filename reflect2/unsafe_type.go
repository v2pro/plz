package reflect2

import (
	"unsafe"
	"reflect"
)

type unsafeType struct {
	reflect.Type
	cfg      *frozenConfig
	rtype    unsafe.Pointer
	ptrRType unsafe.Pointer
}

func newUnsafeType(cfg *frozenConfig, type1 reflect.Type) *unsafeType {
	return &unsafeType{
		Type:     type1,
		cfg:      cfg,
		rtype:    unpackEFace(type1).data,
		ptrRType: unpackEFace(reflect.PtrTo(type1)).data,
	}
}

func (type2 *unsafeType) Type1() reflect.Type {
	return type2.Type
}

func (type2 *unsafeType) UnsafeNew() unsafe.Pointer {
	return unsafe_New(type2.rtype)
}

func (type2 *unsafeType) New() interface{} {
	return packEFace(type2.ptrRType, type2.UnsafeNew())
}

func (type2 *unsafeType) PackEFace(ptr unsafe.Pointer) interface{} {
	return packEFace(type2.rtype, ptr)
}
