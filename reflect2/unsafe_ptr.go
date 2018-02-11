package reflect2

import (
	"unsafe"
	"reflect"
)

type unsafeDirPointerType struct {
	unsafeType
	elemRType unsafe.Pointer
}

type unsafeIndirPointerType struct {
	unsafeType
	elemRType unsafe.Pointer
}

func newUnsafePointerType(type1 reflect.Type) PointerType {
	switch type1.Elem().Kind() {
	case reflect.Ptr, reflect.Map:
		return &unsafeIndirPointerType{
			unsafeType: *newUnsafeType(type1),
			elemRType:  toEface(type1.Elem()).data,
		}
	default:
		return &unsafeDirPointerType{
			unsafeType: *newUnsafeType(type1),
			elemRType:  toEface(type1.Elem()).data,
		}
	}
}

func (type2 *unsafeDirPointerType) Get(obj interface{}) interface{} {
	ptr := type2.UnsafeGet(toEface(obj).data)
	return packEface(type2.elemRType, ptr)
}

func (type2 *unsafeDirPointerType) UnsafeGet(obj unsafe.Pointer) unsafe.Pointer {
	return obj
}

func (type2 *unsafeIndirPointerType) Get(obj interface{}) interface{} {
	ptr := type2.UnsafeGet(toEface(obj).data)
	return packEface(type2.elemRType, ptr)
}

func (type2 *unsafeIndirPointerType) UnsafeGet(obj unsafe.Pointer) unsafe.Pointer {
	return *(*unsafe.Pointer)(obj)
}
