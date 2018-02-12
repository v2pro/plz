package reflect2

import (
	"reflect"
	"unsafe"
)

type unsafeMapType struct {
	unsafeType
	pKeyRType  unsafe.Pointer
	pElemRType unsafe.Pointer
}

func newUnsafeMapType(cfg *frozenConfig, type1 reflect.Type) MapType {
	return &unsafeMapType{
		unsafeType: *newUnsafeType(cfg, type1),
		pKeyRType:  unpackEFace(reflect.PtrTo(type1.Key())).data,
		pElemRType: unpackEFace(reflect.PtrTo(type1.Elem())).data,
	}
}

func (type2 *unsafeMapType) Key() Type {
	return type2.cfg.Type2(type2.Type.Key())
}

func (type2 *unsafeMapType) Elem() Type {
	return type2.cfg.Type2(type2.Type.Elem())
}

func (type2 *unsafeMapType) MakeMap(cap int) interface{} {
	return packEFace(type2.ptrRType, type2.UnsafeMakeMap(cap))
}

func (type2 *unsafeMapType) UnsafeMakeMap(cap int) unsafe.Pointer {
	m := makemap(type2.rtype, cap)
	return unsafe.Pointer(&m)
}

func (type2 *unsafeMapType) Set(obj interface{}, key interface{}, elem interface{}) {
	objEFace := unpackEFace(obj)
	assertType("MapType.Set argument 1", type2.ptrRType, objEFace.rtype)
	keyEFace := unpackEFace(key)
	assertType("MapType.Set argument 2", type2.pKeyRType, keyEFace.rtype)
	elemEFace := unpackEFace(elem)
	assertType("MapType.Set argument 3", type2.pElemRType, elemEFace.rtype)
	type2.UnsafeSet(objEFace.data, keyEFace.data, elemEFace.data)
}

func (type2 *unsafeMapType) UnsafeSet(obj unsafe.Pointer, key unsafe.Pointer, elem unsafe.Pointer) {
	mapassign(type2.rtype, *(*unsafe.Pointer)(obj), key, elem)
}

func (type2 *unsafeMapType) TryGet(obj interface{}, key interface{}) (interface{}, bool) {
	objEFace := unpackEFace(obj)
	assertType("MapType.TryGet argument 1", type2.ptrRType, objEFace.rtype)
	keyEFace := unpackEFace(key)
	assertType("MapType.TryGet argument 2", type2.pKeyRType, keyEFace.rtype)
	elemPtr := type2.UnsafeGet(objEFace.data, keyEFace.data)
	if elemPtr == nil {
		return nil, false
	}
	return packEFace(type2.pElemRType, elemPtr), true
}

func (type2 *unsafeMapType) Get(obj interface{}, key interface{}) interface{} {
	objEFace := unpackEFace(obj)
	assertType("MapType.TryGet argument 1", type2.ptrRType, objEFace.rtype)
	keyEFace := unpackEFace(key)
	assertType("MapType.TryGet argument 2", type2.pKeyRType, keyEFace.rtype)
	elemPtr := type2.UnsafeGet(objEFace.data, keyEFace.data)
	return packEFace(type2.pElemRType, elemPtr)
}

func (type2 *unsafeMapType) UnsafeGet(obj unsafe.Pointer, key unsafe.Pointer) unsafe.Pointer {
	return mapaccess(type2.rtype, *(*unsafe.Pointer)(obj), key)
}

func (type2 *unsafeMapType) Iterate(obj interface{}) MapIterator {
	return type2.UnsafeIterate(unpackEFace(obj).data)
}

func (type2 *unsafeMapType) UnsafeIterate(obj unsafe.Pointer) MapIterator {
	return &unsafeMapIterator{
		hiter:      mapiterinit(type2.rtype, obj),
		pKeyRType:  type2.pKeyRType,
		pElemRType: type2.pElemRType,
	}
}

type unsafeMapIterator struct {
	*hiter
	pKeyRType  unsafe.Pointer
	pElemRType unsafe.Pointer
}

func (iter *unsafeMapIterator) HasNext() bool {
	return iter.key != nil
}

func (iter *unsafeMapIterator) Next() (interface{}, interface{}) {
	key, elem := iter.UnsafeNext()
	return packEFace(iter.pKeyRType, key), packEFace(iter.pElemRType, elem)
}

func (iter *unsafeMapIterator) UnsafeNext() (unsafe.Pointer, unsafe.Pointer) {
	key := iter.key
	elem := iter.value
	mapiternext(iter.hiter)
	return key, elem
}
