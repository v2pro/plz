package reflect2

import (
	"unsafe"
	"reflect"
)

// sliceHeader is a safe version of SliceHeader used within this package.
type sliceHeader struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}

func newUnsafeSliceType(type1 reflect.Type) SliceType {
	sliceType := unsafeSliceType{
		unsafeType: *newUnsafeType(type1),
		elemRType:  toEface(type1.Elem()).data,
		elemSize:   type1.Elem().Size(),
	}
	switch type1.Elem().Kind() {
	case reflect.Interface:
		return &unsafeEfaceSliceType{unsafeSliceType: sliceType}
	default:
		return &sliceType
	}
}

type unsafeSliceType struct {
	unsafeType
	elemRType unsafe.Pointer
	elemSize  uintptr
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

func (type2 *unsafeSliceType) Get(obj interface{}, index int) interface{} {
	elemPtr := type2.UnsafeGet(toEface(obj).data, index)
	return packEface(type2.elemRType, elemPtr)
}

func (type2 *unsafeSliceType) UnsafeGet(obj unsafe.Pointer, index int) unsafe.Pointer {
	header := (*sliceHeader)(obj)
	elemPtr := arrayAt(header.Data, index, type2.elemSize, "i < s.Len")
	return unsafe.Pointer(elemPtr)
}

func (type2 *unsafeSliceType) Append(obj interface{}, elem interface{}) interface{} {
	ptr := type2.UnsafeAppend(toEface(obj).data, toEface(elem).data)
	return packEface(type2.rtype, ptr)
}

func (type2 *unsafeSliceType) UnsafeAppend(obj unsafe.Pointer, elem unsafe.Pointer) unsafe.Pointer {
	header := (*sliceHeader)(obj)
	if header.Cap == header.Len {
		header = type2.grow(header, header.Len + 1)
	}
	type2.UnsafeSet(unsafe.Pointer(header), header.Len, elem)
	header.Len += 1
	return unsafe.Pointer(header)
}

func (type2 *unsafeSliceType) grow(header *sliceHeader, expectedCap int) *sliceHeader {
	newCap := calcNewCap(header.Cap, expectedCap)
	newHeader := (*sliceHeader)(type2.UnsafeMakeSlice(header.Len, newCap))
	typedslicecopy(type2.elemRType, *newHeader, *header)
	return newHeader
}

type unsafeEfaceSliceType struct {
	unsafeSliceType
}

func (type2 *unsafeEfaceSliceType) Set(obj interface{}, index int, elem interface{}) {
	header := (*sliceHeader)(toEface(obj).data)
	elemPtr := arrayAt(header.Data, index, type2.elemSize, "i < s.Len")
	*(*interface{})(elemPtr) = elem
}

func (type2 *unsafeEfaceSliceType) UnsafeSet(obj unsafe.Pointer, index int, elem unsafe.Pointer) {
	header := (*sliceHeader)(obj)
	elemPtr := arrayAt(header.Data, index, type2.elemSize, "i < s.Len")
	(*iface)(elemPtr).data = elem
}

func (type2 *unsafeEfaceSliceType) Get(obj interface{}, index int) interface{} {
	header := (*sliceHeader)(toEface(obj).data)
	elemPtr := arrayAt(header.Data, index, type2.elemSize, "i < s.Len")
	return *(*interface{})(elemPtr)
}

func (type2 *unsafeEfaceSliceType) UnsafeGet(obj unsafe.Pointer, index int) unsafe.Pointer {
	header := (*sliceHeader)(obj)
	elemPtr := arrayAt(header.Data, index, type2.elemSize, "i < s.Len")
	return (*eface)(elemPtr).data
}

func (type2 *unsafeEfaceSliceType) Append(obj interface{}, elem interface{}) interface{} {
	header := (*sliceHeader)(toEface(obj).data)
	if header.Cap == header.Len {
		header = type2.grow(header, header.Len + 1)
	}
	appended := type2.PackEFace(unsafe.Pointer(header))
	type2.Set(appended, header.Len, elem)
	header.Len += 1
	return appended
}

func (type2 *unsafeEfaceSliceType) UnsafeAppend(obj unsafe.Pointer, elem unsafe.Pointer) unsafe.Pointer {
	header := (*sliceHeader)(obj)
	if header.Cap == header.Len {
		header = type2.grow(header, header.Len + 1)
	}
	type2.UnsafeSet(unsafe.Pointer(header), header.Len, elem)
	header.Len += 1
	return unsafe.Pointer(header)
}

func calcNewCap(cap int, expectedCap int) int {
	if cap == 0 {
		cap = expectedCap
	} else {
		for cap < expectedCap {
			if cap < 1024 {
				cap += cap
			} else {
				cap += cap / 4
			}
		}
	}
	return cap
}
