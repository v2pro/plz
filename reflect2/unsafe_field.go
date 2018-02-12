package reflect2

import (
	"reflect"
	"unsafe"
)

type UnsafeStructField struct {
	reflect.StructField
	structType *UnsafeStructType
	rtype      unsafe.Pointer
	ptrRType   unsafe.Pointer
}

func (field *UnsafeStructField) Set(obj interface{}, value interface{}) {
	objEFace := unpackEFace(obj)
	assertType("StructField.Set argument 1", field.structType.ptrRType, objEFace.rtype)
	valueEFace := unpackEFace(value)
	assertType("StructField.Set argument 2", field.ptrRType, valueEFace.rtype)
	field.UnsafeSet(objEFace.data, valueEFace.data)
}

func (field *UnsafeStructField) UnsafeSet(obj unsafe.Pointer, value unsafe.Pointer) {
	fieldPtr := add(obj, field.Offset, "same as non-reflect &v.field")
	typedmemmove(field.rtype, fieldPtr, value)
}

func (field *UnsafeStructField) Get(obj interface{}) interface{} {
	objEFace := unpackEFace(obj)
	assertType("StructField.Get argument 1", field.structType.ptrRType, objEFace.rtype)
	value := field.UnsafeGet(objEFace.data)
	return packEFace(field.ptrRType, value)
}

func (field *UnsafeStructField) UnsafeGet(obj unsafe.Pointer) unsafe.Pointer {
	return add(obj, field.Offset, "same as non-reflect &v.field")
}
