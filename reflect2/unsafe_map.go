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
	case reflect.Ptr:
		return &unsafeIndirKeyMapType{unsafeMapType: mapType}
	case reflect.Interface:
		if type1.Key().NumMethod() == 0 {
			return &unsafeEFaceKeyMapType{unsafeMapType: mapType}
		}
		return &unsafeIFaceKeyMapType{unsafeMapType: mapType}
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

func (type2 *unsafeMapType) UnsafeIterate(obj unsafe.Pointer) MapIterator {
	return &unsafeMapIterator{
		hiter:     mapiterinit(type2.rtype, obj),
		keyRType:  type2.keyRType,
		elemRType: type2.elemRType,
	}
}

type unsafeMapIterator struct {
	*hiter
	keyRType  unsafe.Pointer
	elemRType unsafe.Pointer
}

func (iter *unsafeMapIterator) HasNext() bool {
	return iter.key != nil
}

func (iter *unsafeMapIterator) Next() (interface{}, interface{}) {
	key, elem := iter.UnsafeNext()
	return packEFace(iter.keyRType, key), packEFace(iter.elemRType, elem)
}

func (iter *unsafeMapIterator) UnsafeNext() (unsafe.Pointer, unsafe.Pointer) {
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

func (type2 *unsafeEFaceKeyMapType) Get(obj interface{}, key interface{}) interface{} {
	elemPtr := type2.UnsafeGet(toEFace(obj).data, unsafe.Pointer(&key))
	if elemPtr == nil {
		return nil
	}
	return packEFace(type2.elemRType, elemPtr)
}

func (type2 *unsafeEFaceKeyMapType) Iterate(obj interface{}) MapIterator {
	return type2.UnsafeIterate(toEFace(obj).data)
}

func (type2 *unsafeEFaceKeyMapType) UnsafeIterate(obj unsafe.Pointer) MapIterator {
	return &unsafeEFaceKeyMapIterator{
		unsafeMapIterator{
			hiter:     mapiterinit(type2.rtype, obj),
			keyRType:  type2.keyRType,
			elemRType: type2.elemRType,
		}}
}

type unsafeEFaceKeyMapIterator struct {
	unsafeMapIterator
}

func (iter *unsafeEFaceKeyMapIterator) Next() (interface{}, interface{}) {
	key, elem := iter.UnsafeNext()
	return *(*interface{})(key), packEFace(iter.elemRType, elem)
}

type unsafeIFaceKeyMapType struct {
	unsafeMapType
}

func (type2 *unsafeIFaceKeyMapType) Set(obj interface{}, key interface{}, elem interface{}) {
	keyIFace := unsafe_New(type2.keyRType)
	if key != nil {
		ifaceE2I(type2.keyRType, key, keyIFace)
	}
	type2.UnsafeSet(toEFace(obj).data, keyIFace, toEFace(elem).data)
}

func (type2 *unsafeIFaceKeyMapType) Get(obj interface{}, key interface{}) interface{} {
	keyIFace := unsafe_New(type2.keyRType)
	if key != nil {
		ifaceE2I(type2.keyRType, key, keyIFace)
	}
	elemPtr := type2.UnsafeGet(toEFace(obj).data, keyIFace)
	if elemPtr == nil {
		return nil
	}
	return packEFace(type2.elemRType, elemPtr)
}

func (type2 *unsafeIFaceKeyMapType) Iterate(obj interface{}) MapIterator {
	return type2.UnsafeIterate(toEFace(obj).data)
}

func (type2 *unsafeIFaceKeyMapType) UnsafeIterate(obj unsafe.Pointer) MapIterator {
	return &unsafeIFaceKeyMapIterator{
		unsafeMapIterator{
			hiter:     mapiterinit(type2.rtype, obj),
			keyRType:  type2.keyRType,
			elemRType: type2.elemRType,
		}}
}

type unsafeIFaceKeyMapIterator struct {
	unsafeMapIterator
}

func (iter *unsafeIFaceKeyMapIterator) Next() (interface{}, interface{}) {
	key, elem := iter.UnsafeNext()
	keyIFace := (*iface)(key)
	return packEFace(keyIFace.itab.rtype, keyIFace.data), packEFace(iter.elemRType, elem)
}

type unsafeIndirKeyMapType struct {
	unsafeMapType
}

func (type2 *unsafeIndirKeyMapType) Set(obj interface{}, key interface{}, elem interface{}) {
	type2.UnsafeSet(toEFace(obj).data, unsafe.Pointer(&toEFace(key).data), toEFace(elem).data)
}

func (type2 *unsafeIndirKeyMapType) Get(obj interface{}, key interface{}) interface{} {
	elemPtr := type2.UnsafeGet(toEFace(obj).data, unsafe.Pointer(&toEFace(key).data))
	if elemPtr == nil {
		return nil
	}
	return packEFace(type2.elemRType, elemPtr)
}

func (type2 *unsafeIndirKeyMapType) Iterate(obj interface{}) MapIterator {
	return type2.UnsafeIterate(toEFace(obj).data)
}

func (type2 *unsafeIndirKeyMapType) UnsafeIterate(obj unsafe.Pointer) MapIterator {
	return &unsafeIndirKeyMapIterator{
		unsafeMapIterator{
			hiter:     mapiterinit(type2.rtype, obj),
			keyRType:  type2.keyRType,
			elemRType: type2.elemRType,
		}}
}

type unsafeIndirKeyMapIterator struct {
	unsafeMapIterator
}

func (iter *unsafeIndirKeyMapIterator) Next() (interface{}, interface{}) {
	key, elem := iter.UnsafeNext()
	return packEFace(iter.keyRType, *(*unsafe.Pointer)(key)), packEFace(iter.elemRType, elem)
}