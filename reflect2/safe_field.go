package reflect2

import (
	"reflect"
	"unsafe"
)

type safeField struct {
	reflect.StructField
}

func (field *safeField) Set(obj interface{}, value interface{}) {
	val := reflect.ValueOf(obj).Elem()
	val.FieldByIndex(field.Index).Set(reflect.ValueOf(value).Elem())
}

func (field *safeField) UnsafeSet(obj unsafe.Pointer, value unsafe.Pointer) {
	panic("unsafe operation is not supported")
}

func (field *safeField) Get(obj interface{}) interface{} {
	val := reflect.ValueOf(obj).Elem().FieldByIndex(field.Index)
	ptr := reflect.New(val.Type())
	ptr.Elem().Set(val)
	return ptr.Interface()
}

func (field *safeField) UnsafeGet(obj unsafe.Pointer) unsafe.Pointer {
	panic("does not support unsafe operation")
}