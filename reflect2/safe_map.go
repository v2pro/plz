package reflect2

import (
	"reflect"
	"unsafe"
)

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
	val := reflect.ValueOf(obj)
	val.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(elem))
}

func (type2 *safeMapType) UnsafeSet(obj unsafe.Pointer, key unsafe.Pointer, elem unsafe.Pointer) {
	panic("does not support unsafe operation")
}

func (type2 *safeMapType) Get(obj interface{}, key interface{}) interface{} {
	val := reflect.ValueOf(obj).MapIndex(reflect.ValueOf(key))
	if !val.IsValid() {
		return nil
	}
	return val.Interface()
}

func (type2 *safeMapType) UnsafeGet(obj unsafe.Pointer, key unsafe.Pointer) unsafe.Pointer {
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