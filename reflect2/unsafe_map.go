package reflect2

import (
	"reflect"
	"unsafe"
)

type unsafeMapType struct {
	unsafeType
	keyRType  unsafe.Pointer
	elemRType unsafe.Pointer
}

func newUnsafeMapType(type1 reflect.Type) MapType {
	mapType := unsafeMapType{
		unsafeType: *newUnsafeType(type1),
		keyRType:   toEFace(type1.Key()).data,
		elemRType:  toEFace(type1.Elem()).data,
	}
	switch type1.Key().Kind() {
	case reflect.Interface:
		return &unsafeEFaceKeyMapType{unsafeMapType: mapType}
	}
	return &mapType
}

func (type2 *unsafeMapType) MakeMap(cap int) interface{} {
	return packEFace(type2.rtype, type2.UnsafeMakeMap(cap))
}

func (type2 *unsafeMapType) UnsafeMakeMap(cap int) unsafe.Pointer {
	return makemap(type2.rtype, cap)
}

func (type2 *unsafeMapType) Set(obj interface{}, key interface{}, elem interface{}) {
	type2.UnsafeSet(toEFace(obj).data, toEFace(key).data, toEFace(elem).data)
}

func (type2 *unsafeMapType) UnsafeSet(obj unsafe.Pointer, key unsafe.Pointer, elem unsafe.Pointer) {
	mapassign(type2.rtype, obj, key, elem)
}

func (type2 *unsafeMapType) Get(obj interface{}, key interface{}) interface{} {
	elemPtr := type2.UnsafeGet(toEFace(obj).data, toEFace(key).data)
	if elemPtr == nil {
		return nil
	}
	return packEFace(type2.elemRType, elemPtr)
}

func (type2 *unsafeMapType) UnsafeGet(obj unsafe.Pointer, key unsafe.Pointer) unsafe.Pointer {
	return mapaccess(type2.rtype, obj, key)
}

func (type2 *unsafeMapType) Iterate(obj interface{}) MapIterator {
	return type2.UnsafeIterate(toEFace(obj).data)
}

func (type2 *unsafeMapType) UnsafeIterate(obj unsafe.Pointer) *UnsafeMapIterator {
	return &UnsafeMapIterator{
		hiter:     mapiterinit(type2.rtype, obj),
		keyRType:  type2.keyRType,
		elemRType: type2.elemRType,
	}
}

type UnsafeMapIterator struct {
	*hiter
	keyRType  unsafe.Pointer
	elemRType unsafe.Pointer
}

func (iter *UnsafeMapIterator) HasNext() bool {
	return iter.key != nil
}

func (iter *UnsafeMapIterator) Next() (interface{}, interface{}) {
	key, elem := iter.UnsafeNext()
	return packEFace(iter.keyRType, key), packEFace(iter.elemRType, elem)
}

func (iter *UnsafeMapIterator) UnsafeNext() (unsafe.Pointer, unsafe.Pointer) {
	key := iter.key
	elem := iter.value
	mapiternext(iter.hiter)
	return key, elem
}

type unsafeEFaceKeyMapType struct {
	unsafeMapType
}

func (type2 *unsafeEFaceKeyMapType) Set(obj interface{}, key interface{}, elem interface{}) {
	type2.UnsafeSet(toEFace(obj).data, unsafe.Pointer(&key), toEFace(elem).data)
}

func (type2 *unsafeEFaceKeyMapType) UnsafeSet(obj unsafe.Pointer, key unsafe.Pointer, elem unsafe.Pointer) {
	mapassign(type2.rtype, obj, key, elem)
}

