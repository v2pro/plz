package reflect2

import (
	"reflect"
	"unsafe"
)

type unsafeField struct {
	reflect.StructField
	rtype unsafe.Pointer
}

func (field *unsafeField) Set(obj interface{}, value interface{}) {
	field.UnsafeSet(toEFace(obj).data, toEFace(value).data)
}

func (field *unsafeField) UnsafeSet(obj unsafe.Pointer, value unsafe.Pointer) {
	fieldPtr := add(obj, field.Offset, "same as non-reflect &v.field")
	typedmemmove(field.rtype, fieldPtr, value)
}