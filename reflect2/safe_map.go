package reflect2

import (
	"reflect"
	"unsafe"
)

type safeMapType struct {
	safeType
}

func (type2 *safeMapType) MakeMap(cap int) interface{} {
	ptr := reflect.New(type2.Type)
	ptr.Elem().Set(reflect.MakeMapWithSize(type2.Type, cap))
	return ptr.Interface()
}

func (type2 *safeMapType) UnsafeMakeMap(cap int) unsafe.Pointer {
	panic("does not support unsafe operation")
}

func (type2 *safeMapType) Set(obj interface{}, key interface{}, elem interface{}) {
	keyVal := reflect.ValueOf(key)
	elemVal := reflect.ValueOf(elem)
	val := reflect.ValueOf(obj)
	val.Elem().SetMapIndex(keyVal.Elem(), elemVal.Elem())
}

func (type2 *safeMapType) UnsafeSet(obj unsafe.Pointer, key unsafe.Pointer, elem unsafe.Pointer) {
	panic("does not support unsafe operation")
}

func (type2 *safeMapType) TryGet(obj interface{}, key interface{}) (interface{}, bool) {
	keyVal := reflect.ValueOf(key)
	if key == nil {
		keyVal = reflect.New(type2.Key()).Elem()
	}
	val := reflect.ValueOf(obj).MapIndex(keyVal)
	if !val.IsValid() {
		return nil, false
	}
	return val.Interface(), true
}

func (type2 *safeMapType) Get(obj interface{}, key interface{}) interface{} {
	val := reflect.ValueOf(obj).Elem()
	keyVal := reflect.ValueOf(key).Elem()
	elemVal := val.MapIndex(keyVal)
	if !elemVal.IsValid() {
		ptr := reflect.New(reflect.PtrTo(val.Type().Elem()))
		return ptr.Elem().Interface()
	}
	ptr := reflect.New(elemVal.Type())
	ptr.Elem().Set(elemVal)
	return ptr.Interface()
}

func (type2 *safeMapType) UnsafeGet(obj unsafe.Pointer, key unsafe.Pointer) unsafe.Pointer {
	panic("does not support unsafe operation")
}

func (type2 *safeMapType) Iterate(obj interface{}) MapIterator {
	m := reflect.ValueOf(obj).Elem()
	return &safeMapIterator{
		m:    m,
		keys: m.MapKeys(),
	}
}

func (type2 *safeMapType) UnsafeIterate(obj unsafe.Pointer) MapIterator {
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
	keyPtr := reflect.New(key.Type())
	keyPtr.Elem().Set(key)
	elemPtr := reflect.New(elem.Type())
	elemPtr.Elem().Set(elem)
	return keyPtr.Interface(), elemPtr.Interface()
}

func (iter *safeMapIterator) UnsafeNext() (unsafe.Pointer, unsafe.Pointer) {
	panic("does not support unsafe operation")
}