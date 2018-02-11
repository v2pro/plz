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