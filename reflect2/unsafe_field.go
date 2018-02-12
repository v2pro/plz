package reflect2

import (
	"reflect"
	"unsafe"
)

type unsafeDirField struct {
	reflect.StructField
	rtype unsafe.Pointer
}

func (field *unsafeDirField) Set(obj interface{}, value interface{}) {
	field.UnsafeSet(unpackEFace(obj).data, unpackEFace(value).data)
}

func (field *unsafeDirField) UnsafeSet(obj unsafe.Pointer, value unsafe.Pointer) {
	fieldPtr := add(obj, field.Offset, "same as non-reflect &v.field")
	typedmemmove(field.rtype, fieldPtr, value)
}

func (field *unsafeDirField) Get(obj interface{}) interface{} {
	value := field.UnsafeGet(unpackEFace(obj).data)
	return packEFace(field.rtype, value)
}

func (field *unsafeDirField) UnsafeGet(obj unsafe.Pointer) unsafe.Pointer {
	fieldPtr := add(obj, field.Offset, "same as non-reflect &v.field")
	return fieldPtr
}

type unsafeEFaceField struct {
	reflect.StructField
}

func (field *unsafeEFaceField) Set(obj interface{}, value interface{}) {
	field.UnsafeSet(unpackEFace(obj).data, unsafe.Pointer(&value))
}

func (field *unsafeEFaceField) UnsafeSet(obj unsafe.Pointer, value unsafe.Pointer) {
	fieldPtr := add(obj, field.Offset, "same as non-reflect &v.field")
	*(*interface{})(fieldPtr) = *(*interface{})(value)
}

func (field *unsafeEFaceField) Get(obj interface{}) interface{} {
	value := field.UnsafeGet(unpackEFace(obj).data)
	return *(*interface{})(value)
}

func (field *unsafeEFaceField) UnsafeGet(obj unsafe.Pointer) unsafe.Pointer {
	fieldPtr := add(obj, field.Offset, "same as non-reflect &v.field")
	return fieldPtr
}