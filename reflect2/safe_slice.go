package reflect2

import (
	"reflect"
	"unsafe"
)

type safeSliceType struct {
	safeType
}

func (type2 *safeSliceType) Set(obj interface{}, index int, value interface{}) {
	val := reflect.ValueOf(obj).Elem()
	elem := reflect.ValueOf(value).Elem()
	val.Index(index).Set(elem)
}

func (type2 *safeSliceType) UnsafeSet(obj unsafe.Pointer, index int, value unsafe.Pointer) {
	panic("does not support unsafe operation")
}

func (type2 *safeSliceType) Get(obj interface{}, index int) interface{} {
	val := reflect.ValueOf(obj).Elem()
	elem := val.Index(index)
	ptr := reflect.New(elem.Type())
	ptr.Elem().Set(elem)
	return ptr.Interface()
}

func (type2 *safeSliceType) UnsafeGet(obj unsafe.Pointer, index int) unsafe.Pointer {
	panic("does not support unsafe operation")
}

func (type2 *safeSliceType) MakeSlice(length int, cap int) interface{} {
	val := reflect.MakeSlice(type2.Type, length, cap)
	ptr := reflect.New(val.Type())
	ptr.Elem().Set(val)
	return ptr.Interface()
}

func (type2 *safeSliceType) UnsafeMakeSlice(length int, cap int) unsafe.Pointer {
	panic("does not support unsafe operation")
}

func (type2 *safeSliceType) Append(obj interface{}, elem interface{}) interface{} {
	val := reflect.ValueOf(obj).Elem()
	elemVal := reflect.ValueOf(elem).Elem()
	val = reflect.Append(val, elemVal)
	ptr := reflect.New(val.Type())
	ptr.Elem().Set(val)
	return ptr.Interface()
}

func (type2 *safeSliceType) UnsafeAppend(obj unsafe.Pointer, elem unsafe.Pointer) unsafe.Pointer {
	panic("does not support unsafe operation")
}

func (type2 *safeSliceType) SetNil(obj interface{}) {
	val := reflect.ValueOf(obj).Elem()
	val.Set(reflect.Zero(val.Type()))
}

func (type2 *safeSliceType) UnsafeSetNil(ptr unsafe.Pointer) {
	panic("does not support unsafe operation")
}

func (type2 *safeSliceType) LengthOf(obj interface{}) int {
	return reflect.ValueOf(obj).Elem().Len()
}

func (type2 *safeSliceType) UnsafeLengthOf(ptr unsafe.Pointer) int {
	panic("does not support unsafe operation")
}