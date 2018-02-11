package reflect2

import (
	"reflect"
	"unsafe"
)

type unsafeMapType struct {
	unsafeType
}

func newUnsafeMapType(type1 reflect.Type) *unsafeMapType {
	return &unsafeMapType{
		unsafeType: *newUnsafeType(type1),
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