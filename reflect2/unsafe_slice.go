package reflect2

import (
	"unsafe"
	"reflect"
)

type unsafeSliceType struct {
	unsafeType
	elemRType unsafe.Pointer
	elemSize  uintptr
}

// sliceHeader is a safe version of SliceHeader used within this package.
type sliceHeader struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}

func newUnsafeSliceType(type1 reflect.Type) *unsafeSliceType {
	return &unsafeSliceType{
		unsafeType: *newUnsafeType(type1),
		elemRType:  toEface(type1.Elem()).data,
		elemSize:   type1.Elem().Size(),
	}
}

func (type2 *unsafeSliceType) MakeSlice(length int, cap int) interface{} {
	return packEface(type2.rtype, type2.UnsafeMakeSlice(length, cap))
}

func (type2 *unsafeSliceType) UnsafeMakeSlice(length int, cap int) unsafe.Pointer {
	header := &sliceHeader{unsafe_NewArray(type2.elemRType, cap), length, cap}
	return unsafe.Pointer(header)
}

func (type2 *unsafeSliceType) Set(obj interface{}, index int, elem interface{}) {
	type2.UnsafeSet(toEface(obj).data, index, toEface(elem).data)
}

func (type2 *unsafeSliceType) UnsafeSet(obj unsafe.Pointer, index int, elem unsafe.Pointer) {
	header := (*sliceHeader)(obj)
	elemPtr := arrayAt(header.Data, index, type2.elemSize, "i < s.Len")
	typedmemmove(type2.elemRType, elemPtr, elem)
}
