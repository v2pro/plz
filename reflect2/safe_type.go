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

type safeMapType struct {
	safeType
}

func (type2 *safeMapType) MakeMap(cap int) interface{} {
	return reflect.MakeMapWithSize(type2.Type, cap).Interface()
}

func (type2 *safeMapType) UnsafeMakeMap(cap int) unsafe.Pointer {
	panic("does not support unsafe operation")
}

func (type2 *safeMapType) Set(obj interface{}, key interface{}, elem interface{}) {
	reflect.ValueOf(obj).SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(elem))
}

func (type2 *safeMapType) UnsafeSet(obj unsafe.Pointer, key unsafe.Pointer, elem unsafe.Pointer) {
	panic("does not support unsafe operation")
}


func (type2 *safeMapType) Iterate(obj interface{}) MapIterator {
	m := reflect.ValueOf(obj)
	return &safeMapIterator{
		m:    m,
		keys: m.MapKeys(),
	}
}

func (type2 *safeMapType) UnsafeIterate(obj unsafe.Pointer) *UnsafeMapIterator {
	panic("does not support unsafe operation")
}

type safeMapIterator struct {
	i int
	m reflect.Value
	keys []reflect.Value
}

func (iter *safeMapIterator) HasNext() bool {
	return iter.i != len(iter.keys)
}

func (iter *safeMapIterator) Next() (interface{}, interface{}) {
	key := iter.keys[iter.i]
	elem := iter.m.MapIndex(key)
	iter.i += 1
	return key.Interface(), elem.Interface()
}