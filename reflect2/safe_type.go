package reflect2

import (
	"reflect"
	"unsafe"
)

type safeType struct {
	reflect.Type
}

func (type2 *safeType) New() interface{} {
	return reflect.New(type2.Type).Interface()
}

func (type2 *safeType) UnsafeNew() unsafe.Pointer {
	panic("does not support unsafe operation")
}

func (type2 *safeType) Type1() reflect.Type {
	return type2.Type
}

func (type2 *safeType) PackEFace(ptr unsafe.Pointer) interface{} {
	panic("does not support unsafe operation")
}

func (type2 *safeType) FieldByName(name string) StructField {
	field, found := type2.Type.FieldByName(name)
	if !found {
		panic("field " + name + " not found")
	}
	return &safeField{StructField: field}
}

func (type2 *safeType) Set(obj interface{}, index int, value interface{}) {
	reflect.ValueOf(obj).Elem().Index(index).Set(reflect.ValueOf(value))
}

func (type2 *safeType) UnsafeSet(obj unsafe.Pointer, index int, value unsafe.Pointer) {
	panic("does not support unsafe operation")
}

func (type2 *safeType) Get(obj interface{}, index int) interface{} {
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	return val.Index(index).Interface()
}

func (type2 *safeType) UnsafeGet(obj unsafe.Pointer, index int) unsafe.Pointer {
	panic("does not support unsafe operation")
}

func (type2 *safeType) MakeSlice(length int, cap int) interface{} {
	return reflect.MakeSlice(type2.Type, length, cap).Interface()
}

func (type2 *safeType) UnsafeMakeSlice(length int, cap int) unsafe.Pointer {
	panic("does not support unsafe operation")
}

func (type2 *safeType) Append(obj interface{}, elem interface{}) interface{} {
	return reflect.Append(reflect.ValueOf(obj), reflect.ValueOf(elem)).Interface()
}

func (type2 *safeType) UnsafeAppend(obj unsafe.Pointer, elem unsafe.Pointer) unsafe.Pointer{
	panic("does not support unsafe operation")
}