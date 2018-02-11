package reflect2

import (
	"reflect"
	"unsafe"
)

type safeType struct {
	reflect.Type
}

func (type2 safeType) New() interface{} {
	return reflect.New(type2.Type).Interface()
}

func (type2 safeType) UnsafeNew() unsafe.Pointer {
	panic("does not support unsafe operation")
}

func (type2 safeType) Type1() reflect.Type {
	return type2.Type
}

func (type2 safeType) FieldByName(fieldName string) StructField {
	panic("not implemented yet")
}
