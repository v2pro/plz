package reflect2

import (
	"unsafe"
	"reflect"
)

type UnsafePtrType struct {
	unsafeType
}

func newUnsafePtrType(cfg *frozenConfig, type1 reflect.Type) *UnsafePtrType {
	return &UnsafePtrType{
		unsafeType: *newUnsafeType(cfg, type1),
	}
}

func (type2 *UnsafePtrType) Indirect(obj interface{}) interface{} {
	objEFace := unpackEFace(obj)
	assertType("Type.Indirect argument 1", type2.ptrRType, objEFace.rtype)
	return type2.UnsafeIndirect(objEFace.data)
}

func (type2 *UnsafePtrType) UnsafeIndirect(ptr unsafe.Pointer) interface{} {
	return packEFace(type2.rtype, *(*unsafe.Pointer)(ptr))
}
