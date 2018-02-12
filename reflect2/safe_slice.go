package reflect2

import (
	"reflect"
	"unsafe"
)

type safeSliceType struct {
	safeType
}

func (type2 *safeSliceType) Set(obj interface{}, index int, value interface{}) {
	reflect.ValueOf(obj).Index(index).Set(reflect.ValueOf(value).Elem())
}

func (type2 *safeSliceType) UnsafeSet(obj unsafe.Pointer, index int, value unsafe.Pointer) {
	panic("does not support unsafe operation")
}

func (type2 *safeSliceType) Get(obj interface{}, index int) interface{} {
	val := reflect.ValueOf(obj).Index(index)
	v := val.Interface()
	return &v
}

func (type2 *safeSliceType) UnsafeGet(obj unsafe.Pointer, index int) unsafe.Pointer {
	panic("does not support unsafe operation")
}

func (type2 *safeSliceType) MakeSlice(length int, cap int) interface{} {
	return reflect.MakeSlice(type2.Type, length, cap).Interface()
}

func (type2 *safeSliceType) UnsafeMakeSlice(length int, cap int) unsafe.Pointer {
	panic("does not support unsafe operation")
}

func (type2 *safeSliceType) Append(obj interface{}, elem interface{}) interface{} {
	return reflect.Append(reflect.ValueOf(obj), reflect.ValueOf(elem)).Interface()
}

func (type2 *safeSliceType) UnsafeAppend(obj unsafe.Pointer, elem unsafe.Pointer) unsafe.Pointer{
	panic("does not support unsafe operation")
}