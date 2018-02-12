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
			elemRType:  toEFace(type1.Elem()).data,
		}
	case reflect.Interface:
		if type1.Elem().NumMethod() == 0 {
			return &unsafeEFacePointerType{
				unsafeType: *newUnsafeType(type1),
			}
		}
		return &unsafeIFacePointerType{
			unsafeType: *newUnsafeType(type1),
		}
	default:
		return &unsafeDirPointerType{
			unsafeType: *newUnsafeType(type1),
			elemRType:  toEFace(type1.Elem()).data,
		}
	}
}

type unsafeDirPointerType struct {
	unsafeType
	elemRType unsafe.Pointer
}

func (type2 *unsafeDirPointerType) Get(obj interface{}) interface{} {
	ptr := type2.UnsafeGet(toEFace(obj).data)
	return packEFace(type2.elemRType, ptr)
}

func (type2 *unsafeDirPointerType) UnsafeGet(obj unsafe.Pointer) unsafe.Pointer {
	return obj
}

type unsafeIndirPointerType struct {
	unsafeType
	elemRType unsafe.Pointer
}

func (type2 *unsafeIndirPointerType) Get(obj interface{}) interface{} {
	ptr := type2.UnsafeGet(toEFace(obj).data)
	return packEFace(type2.elemRType, ptr)
}

func (type2 *unsafeIndirPointerType) UnsafeGet(obj unsafe.Pointer) unsafe.Pointer {
	return *(*unsafe.Pointer)(obj)
}

type unsafeEFacePointerType struct {
	unsafeType
}

func (type2 *unsafeEFacePointerType) Get(obj interface{}) interface{} {
	ptr := (*interface{})(toEFace(obj).data)
	return *ptr
}

func (type2 *unsafeEFacePointerType) UnsafeGet(obj unsafe.Pointer) unsafe.Pointer {
	return (*eface)(obj).data
}

type unsafeIFacePointerType struct {
	unsafeType
}

func (type2 *unsafeIFacePointerType) Get(obj interface{}) interface{} {
	ptr := (*iface)(toEFace(obj).data)
	return packEFace(ptr.itab.rtype, ptr.data)
}

func (type2 *unsafeIFacePointerType) UnsafeGet(obj unsafe.Pointer) unsafe.Pointer {
	return (*iface)(obj).data
}
