package reflect2

import (
	"reflect"
	"unsafe"
)

type unsafeMapType struct {
	unsafeType
	keyRType      unsafe.Pointer
	elemRType     unsafe.Pointer
	keyEmbedType  embedType
	elemEmbedType embedType
}

type embedType interface {
	Pack(ptr unsafe.Pointer) interface{}
	Unpack(obj interface{}) unsafe.Pointer
}

type dirEmbedType struct {
	rtype unsafe.Pointer
}

func (embedType *dirEmbedType) Pack(ptr unsafe.Pointer) interface{} {
	return packEFace(embedType.rtype, ptr)
}

func (embedType *dirEmbedType) Unpack(obj interface{}) unsafe.Pointer {
	if obj == nil {
		var nilPtr unsafe.Pointer
		return unsafe.Pointer(&nilPtr)
	}
	return unpackEFace(obj).data
}

type indirEmbedType struct {
	rtype unsafe.Pointer
}

func (embedType *indirEmbedType) Pack(ptr unsafe.Pointer) interface{} {
	return packEFace(embedType.rtype, *(*unsafe.Pointer)(ptr))
}

func (embedType *indirEmbedType) Unpack(obj interface{}) unsafe.Pointer {
	return unsafe.Pointer(&unpackEFace(obj).data)
}

type efaceEmbedType struct {
	rtype unsafe.Pointer
}

func (embedType *efaceEmbedType) Pack(ptr unsafe.Pointer) interface{} {
	return *(*interface{})(ptr)
}

func (embedType *efaceEmbedType) Unpack(obj interface{}) unsafe.Pointer {
	return unsafe.Pointer(&obj)
}

type ifaceEmbedType struct {
	rtype unsafe.Pointer
}

func (embedType *ifaceEmbedType) Pack(ptr unsafe.Pointer) interface{} {
	iface := (*iface)(ptr)
	return packEFace(iface.itab.rtype, iface.data)
}

func (embedType *ifaceEmbedType) Unpack(obj interface{}) unsafe.Pointer {
	iface := unsafe_New(embedType.rtype)
	if obj != nil {
		ifaceE2I(embedType.rtype, obj, iface)
	}
	return iface
}

func newEmbedType(type1 reflect.Type) embedType {
	rtype := unpackEFace(type1).data
	switch type1.Kind() {
	case reflect.Map, reflect.Ptr, reflect.Chan, reflect.Func:
		return &indirEmbedType{rtype: rtype}
	case reflect.Interface:
		if type1.NumMethod() == 0 {
			return &efaceEmbedType{rtype: rtype}
		}
		return &ifaceEmbedType{rtype: rtype}
	case reflect.Struct:
		if type1.NumField() == 1 {
			firstFieldKind := type1.Field(0).Type.Kind()
			if firstFieldKind == reflect.Ptr || firstFieldKind == reflect.Map ||
				firstFieldKind == reflect.Chan || firstFieldKind == reflect.Func {
				return &indirEmbedType{rtype: rtype}
			}
		}
		return &dirEmbedType{rtype: rtype}
	case reflect.Array:
		if type1.Len() == 1 {
			elemKind := type1.Elem().Kind()
			if elemKind == reflect.Ptr || elemKind == reflect.Map ||
				elemKind == reflect.Chan || elemKind == reflect.Func {
				return &indirEmbedType{rtype: rtype}
			}
		}
		return &dirEmbedType{rtype: rtype}
	default:
		return &dirEmbedType{rtype: rtype}
	}
}

func newUnsafeMapType(type1 reflect.Type) MapType {
	return &unsafeMapType{
		unsafeType:    *newUnsafeType(type1),
		keyRType:      unpackEFace(type1.Key()).data,
		elemRType:     unpackEFace(type1.Elem()).data,
		keyEmbedType:  newEmbedType(type1.Key()),
		elemEmbedType: newEmbedType(type1.Elem()),
	}
}

func (type2 *unsafeMapType) MakeMap(cap int) interface{} {
	return packEFace(type2.rtype, makemap(type2.rtype, cap))
}

func (type2 *unsafeMapType) UnsafeMakeMap(cap int) unsafe.Pointer {
	m := makemap(type2.rtype, cap)
	return unsafe.Pointer(&m)
}

func (type2 *unsafeMapType) PackEFace(ptr unsafe.Pointer) interface{} {
	return type2.unsafeType.PackEFace(*(*unsafe.Pointer)(ptr))
}

func (type2 *unsafeMapType) Set(obj interface{}, key interface{}, elem interface{}) {
	var elemPtr unsafe.Pointer
		elemPtr = type2.elemEmbedType.Unpack(elem)
	elemPtr = type2.elemEmbedType.Unpack(elem)
	mapassign(type2.rtype, unpackEFace(obj).data,
		type2.keyEmbedType.Unpack(key),
		elemPtr)
}

func (type2 *unsafeMapType) UnsafeSet(obj unsafe.Pointer, key unsafe.Pointer, elem unsafe.Pointer) {
	mapassign(type2.rtype, *(*unsafe.Pointer)(obj), key, elem)
}

func (type2 *unsafeMapType) TryGet(obj interface{}, key interface{}) (interface{}, bool) {
	elemPtr := mapaccess(type2.rtype, unpackEFace(obj).data, type2.keyEmbedType.Unpack(key))
	if elemPtr == nil {
		return nil, false
	}
	return type2.elemEmbedType.Pack(elemPtr), true
}

func (type2 *unsafeMapType) Get(obj interface{}, key interface{}) interface{} {
	elemPtr := mapaccess(type2.rtype, unpackEFace(obj).data, type2.keyEmbedType.Unpack(key))
	if elemPtr == nil {
		return nil
	}
	return type2.elemEmbedType.Pack(elemPtr)
}

func (type2 *unsafeMapType) UnsafeGet(obj unsafe.Pointer, key unsafe.Pointer) unsafe.Pointer {
	return mapaccess(type2.rtype, *(*unsafe.Pointer)(obj), key)
}

func (type2 *unsafeMapType) Iterate(obj interface{}) MapIterator {
	return &unsafeMapIterator{
		hiter:         mapiterinit(type2.rtype, unpackEFace(obj).data),
		keyEmbedType:  type2.keyEmbedType,
		elemEmbedType: type2.elemEmbedType,
	}
}

func (type2 *unsafeMapType) UnsafeIterate(obj unsafe.Pointer) MapIterator {
	return &unsafeMapIterator{
		hiter:         mapiterinit(type2.rtype, *(*unsafe.Pointer)(obj)),
		keyEmbedType:  type2.keyEmbedType,
		elemEmbedType: type2.elemEmbedType,
	}
}

type unsafeMapIterator struct {
	*hiter
	keyEmbedType  embedType
	elemEmbedType embedType
}

func (iter *unsafeMapIterator) HasNext() bool {
	return iter.key != nil
}

func (iter *unsafeMapIterator) Next() (interface{}, interface{}) {
	key, elem := iter.UnsafeNext()
	return iter.keyEmbedType.Pack(key), iter.elemEmbedType.Pack(elem)
}

func (iter *unsafeMapIterator) UnsafeNext() (unsafe.Pointer, unsafe.Pointer) {
	key := iter.key
	elem := iter.value
	mapiternext(iter.hiter)
	return key, elem
}
