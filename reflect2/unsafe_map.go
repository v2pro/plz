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

func newUnsafeMapType(type1 reflect.Type) *unsafeMapType {
	return &unsafeMapType{
		unsafeType: *newUnsafeType(type1),
		keyRType:   toEface(type1.Key()).data,
		elemRType:  toEface(type1.Elem()).data,
	}
}

func (type2 *unsafeMapType) MakeMap(cap int) interface{} {
	return packEface(type2.rtype, type2.UnsafeMakeMap(cap))
}

func (type2 *unsafeMapType) UnsafeMakeMap(cap int) unsafe.Pointer {
	return makemap(type2.rtype, cap)
}

func (type2 *unsafeMapType) Set(obj interface{}, key interface{}, elem interface{}) {
	type2.UnsafeSet(toEface(obj).data, toEface(key).data, toEface(elem).data)
}

func (type2 *unsafeMapType) UnsafeSet(obj unsafe.Pointer, key unsafe.Pointer, elem unsafe.Pointer) {
	mapassign(type2.rtype, obj, key, elem)
}

func (type2 *unsafeMapType) Iterate(obj interface{}) MapIterator {
	return type2.UnsafeIterate(toEface(obj).data)
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
	return packEface(iter.keyRType, key), packEface(iter.elemRType, elem)
}

func (iter *UnsafeMapIterator) UnsafeNext() (unsafe.Pointer, unsafe.Pointer) {
	key := iter.key
	elem := iter.value
	mapiternext(iter.hiter)
	return key, elem
}
