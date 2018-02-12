package reflect2

import (
	"unsafe"
	"reflect"
)

func newUnsafePointerType(type1 reflect.Type) PointerType {
	switch type1.Elem().Kind() {
	case reflect.Ptr, reflect.Map:
		return &unsafeIndirPointerType{
			unsafeType: *newUnsafeType(type1),
			elemRType:  toEface(type1.Elem()).data,
		}
	case reflect.Interface:
		if type1.Elem().NumMethod() == 0 {
			return &unsafeEfacePointerType{
				unsafeType: *newUnsafeType(type1),
			}
		}
		return &unsafeIfacePointerType{
			unsafeType: *newUnsafeType(type1),
		}
	default:
		return &unsafeDirPointerType{
			unsafeType: *newUnsafeType(type1),
			elemRType:  toEface(type1.Elem()).data,
		}
	}
}

type unsafeDirPointerType struct {
	unsafeType
	elemRType unsafe.Pointer
}

func (type2 *unsafeDirPointerType) Get(obj interface{}) interface{} {
	ptr := type2.UnsafeGet(toEface(obj).data)
	return packEface(type2.elemRType, ptr)
}

func (type2 *unsafeDirPointerType) UnsafeGet(obj unsafe.Pointer) unsafe.Pointer {
	return obj
}

type unsafeIndirPointerType struct {
	unsafeType
	elemRType unsafe.Pointer
}

func (type2 *unsafeIndirPointerType) Get(obj interface{}) interface{} {
	ptr := type2.UnsafeGet(toEface(obj).data)
	return packEface(type2.elemRType, ptr)
}

func (type2 *unsafeIndirPointerType) UnsafeGet(obj unsafe.Pointer) unsafe.Pointer {
	return *(*unsafe.Pointer)(obj)
}

type unsafeEfacePointerType struct {
	unsafeType
}

func (type2 *unsafeEfacePointerType) Get(obj interface{}) interface{} {
	ptr := (*interface{})(toEface(obj).data)
	return *ptr
}

func (type2 *unsafeEfacePointerType) UnsafeGet(obj unsafe.Pointer) unsafe.Pointer {
	return (*eface)(obj).data
}

type unsafeIfacePointerType struct {
	unsafeType
}

func (type2 *unsafeIfacePointerType) Get(obj interface{}) interface{} {
	ptr := (*iface)(toEface(obj).data)
	return packEface(ptr.itab.rtype, ptr.data)
}

func (type2 *unsafeIfacePointerType) UnsafeGet(obj unsafe.Pointer) unsafe.Pointer {
	return (*iface)(obj).data
}
