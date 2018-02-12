package reflect2

import (
	"unsafe"
	"reflect"
)

type unsafeArrayType struct {
	unsafeType
	elemRType unsafe.Pointer
	pElemRType unsafe.Pointer
	elemSize uintptr
}

func newUnsafeArrayType(cfg *frozenConfig, type1 reflect.Type) *unsafeArrayType {
	return &unsafeArrayType{
		unsafeType: *newUnsafeType(cfg, type1),
		elemRType:  unpackEFace(type1.Elem()).data,
		pElemRType:  unpackEFace(reflect.PtrTo(type1.Elem())).data,
		elemSize:   type1.Elem().Size(),
	}
}

func (type2 *unsafeArrayType) Set(obj interface{}, index int, elem interface{}) {
	objEFace := unpackEFace(obj)
	assertType("ArrayType.Set argument 1", type2.ptrRType, objEFace.rtype)
	elemEFace := unpackEFace(elem)
	assertType("ArrayType.Set argument 3", type2.pElemRType, elemEFace.rtype)
	type2.UnsafeSet(objEFace.data, index, elemEFace.data)
}

func (type2 *unsafeArrayType) UnsafeSet(obj unsafe.Pointer, index int, elem unsafe.Pointer) {
	elemPtr := arrayAt(obj, index, type2.elemSize, "i < s.Len")
	typedmemmove(type2.elemRType, elemPtr, elem)
}

func (type2 *unsafeArrayType) Get(obj interface{}, index int) interface{} {
	objEFace := unpackEFace(obj)
	assertType("ArrayType.Set argument 1", type2.ptrRType, objEFace.rtype)
	elemPtr := type2.UnsafeGet(objEFace.data, index)
	return packEFace(type2.pElemRType, elemPtr)
}

func (type2 *unsafeArrayType) UnsafeGet(obj unsafe.Pointer, index int) unsafe.Pointer {
	return arrayAt(obj, index, type2.elemSize, "i < s.Len")
}