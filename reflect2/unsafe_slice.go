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
		elemRType:  toEFace(type1.Elem()).data,
		elemSize:   type1.Elem().Size(),
	}
	switch type1.Elem().Kind() {
	case reflect.Map:
		return &unsafeIndirSliceType{unsafeSliceType: sliceType}
	case reflect.Interface:
		if type1.Elem().NumMethod() == 0 {
			return &unsafeEFaceSliceType{unsafeSliceType: sliceType}
		}
		return &unsafeIFaceSliceType{unsafeSliceType: sliceType}
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
	return packEFace(type2.rtype, type2.UnsafeMakeSlice(length, cap))
}

func (type2 *unsafeSliceType) UnsafeMakeSlice(length int, cap int) unsafe.Pointer {
	header := &sliceHeader{unsafe_NewArray(type2.elemRType, cap), length, cap}
	return unsafe.Pointer(header)
}

func (type2 *unsafeSliceType) Set(obj interface{}, index int, elem interface{}) {
	type2.UnsafeSet(toEFace(obj).data, index, toEFace(elem).data)
}

func (type2 *unsafeSliceType) UnsafeSet(obj unsafe.Pointer, index int, elem unsafe.Pointer) {
	header := (*sliceHeader)(obj)
	elemPtr := arrayAt(header.Data, index, type2.elemSize, "i < s.Len")
	typedmemmove(type2.elemRType, elemPtr, elem)
}

func (type2 *unsafeSliceType) Get(obj interface{}, index int) interface{} {
	elemPtr := type2.UnsafeGet(toEFace(obj).data, index)
	return packEFace(type2.elemRType, elemPtr)
}

func (type2 *unsafeSliceType) UnsafeGet(obj unsafe.Pointer, index int) unsafe.Pointer {
	header := (*sliceHeader)(obj)
	elemPtr := arrayAt(header.Data, index, type2.elemSize, "i < s.Len")
	return unsafe.Pointer(elemPtr)
}

func (type2 *unsafeSliceType) Append(obj interface{}, elem interface{}) interface{} {
	ptr := type2.UnsafeAppend(toEFace(obj).data, toEFace(elem).data)
	return packEFace(type2.rtype, ptr)
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

type unsafeEFaceSliceType struct {
	unsafeSliceType
}

func (type2 *unsafeEFaceSliceType) Set(obj interface{}, index int, elem interface{}) {
	header := (*sliceHeader)(toEFace(obj).data)
	elemPtr := arrayAt(header.Data, index, type2.elemSize, "i < s.Len")
	*(*interface{})(elemPtr) = elem
}

func (type2 *unsafeEFaceSliceType) UnsafeSet(obj unsafe.Pointer, index int, elem unsafe.Pointer) {
	header := (*sliceHeader)(obj)
	elemPtr := arrayAt(header.Data, index, type2.elemSize, "i < s.Len")
	(*eface)(elemPtr).data = elem
}

func (type2 *unsafeEFaceSliceType) Get(obj interface{}, index int) interface{} {
	header := (*sliceHeader)(toEFace(obj).data)
	elemPtr := arrayAt(header.Data, index, type2.elemSize, "i < s.Len")
	return *(*interface{})(elemPtr)
}

func (type2 *unsafeEFaceSliceType) UnsafeGet(obj unsafe.Pointer, index int) unsafe.Pointer {
	header := (*sliceHeader)(obj)
	elemPtr := arrayAt(header.Data, index, type2.elemSize, "i < s.Len")
	return (*eface)(elemPtr).data
}

func (type2 *unsafeEFaceSliceType) Append(obj interface{}, elem interface{}) interface{} {
	header := (*sliceHeader)(toEFace(obj).data)
	if header.Cap == header.Len {
		header = type2.grow(header, header.Len + 1)
	}
	appended := type2.PackEFace(unsafe.Pointer(header))
	type2.Set(appended, header.Len, elem)
	header.Len += 1
	return appended
}

func (type2 *unsafeEFaceSliceType) UnsafeAppend(obj unsafe.Pointer, elem unsafe.Pointer) unsafe.Pointer {
	header := (*sliceHeader)(obj)
	if header.Cap == header.Len {
		header = type2.grow(header, header.Len + 1)
	}
	type2.UnsafeSet(unsafe.Pointer(header), header.Len, elem)
	header.Len += 1
	return unsafe.Pointer(header)
}

type unsafeIFaceSliceType struct {
	unsafeSliceType
}

func (type2 *unsafeIFaceSliceType) Set(obj interface{}, index int, elem interface{}) {
	header := (*sliceHeader)(toEFace(obj).data)
	elemPtr := arrayAt(header.Data, index, type2.elemSize, "i < s.Len")
	ifaceE2I(type2.elemRType, elem, elemPtr)
}

func (type2 *unsafeIFaceSliceType) UnsafeSet(obj unsafe.Pointer, index int, elem unsafe.Pointer) {
	header := (*sliceHeader)(obj)
	elemPtr := arrayAt(header.Data, index, type2.elemSize, "i < s.Len")
	(*iface)(elemPtr).data = elem
}

func (type2 *unsafeIFaceSliceType) Get(obj interface{}, index int) interface{} {
	header := (*sliceHeader)(toEFace(obj).data)
	elemPtr := arrayAt(header.Data, index, type2.elemSize, "i < s.Len")
	elemIFace := (*iface)(elemPtr)
	if elemIFace.data == nil {
		return nil
	}
	return packEFace(elemIFace.itab.rtype, elemIFace.data)
}

func (type2 *unsafeIFaceSliceType) UnsafeGet(obj unsafe.Pointer, index int) unsafe.Pointer {
	header := (*sliceHeader)(obj)
	elemPtr := arrayAt(header.Data, index, type2.elemSize, "i < s.Len")
	return (*iface)(elemPtr).data
}

func (type2 *unsafeIFaceSliceType) Append(obj interface{}, elem interface{}) interface{} {
	header := (*sliceHeader)(toEFace(obj).data)
	if header.Cap == header.Len {
		header = type2.grow(header, header.Len + 1)
	}
	appended := type2.PackEFace(unsafe.Pointer(header))
	type2.Set(appended, header.Len, elem)
	header.Len += 1
	return appended
}

func (type2 *unsafeIFaceSliceType) UnsafeAppend(obj unsafe.Pointer, elem unsafe.Pointer) unsafe.Pointer {
	header := (*sliceHeader)(obj)
	if header.Cap == header.Len {
		header = type2.grow(header, header.Len + 1)
	}
	type2.UnsafeSet(unsafe.Pointer(header), header.Len, elem)
	header.Len += 1
	return unsafe.Pointer(header)
}

type unsafeIndirSliceType struct {
	unsafeSliceType
}

func (type2 *unsafeIndirSliceType) Set(obj interface{}, index int, elem interface{}) {
	type2.UnsafeSet(toEFace(obj).data, index, toEFace(elem).data)
}

func (type2 *unsafeIndirSliceType) UnsafeSet(obj unsafe.Pointer, index int, elem unsafe.Pointer) {
	header := (*sliceHeader)(obj)
	elemPtr := arrayAt(header.Data, index, type2.elemSize, "i < s.Len")
	*(*unsafe.Pointer)(elemPtr) = elem
}

func (type2 *unsafeIndirSliceType) Get(obj interface{}, index int) interface{} {
	elemPtr := type2.UnsafeGet(toEFace(obj).data, index)
	return packEFace(type2.elemRType, elemPtr)
}

func (type2 *unsafeIndirSliceType) UnsafeGet(obj unsafe.Pointer, index int) unsafe.Pointer {
	header := (*sliceHeader)(obj)
	elemPtr := arrayAt(header.Data, index, type2.elemSize, "i < s.Len")
	return *(*unsafe.Pointer)(elemPtr)
}

func (type2 *unsafeIndirSliceType) Append(obj interface{}, elem interface{}) interface{} {
	ptr := type2.UnsafeAppend(toEFace(obj).data, toEFace(elem).data)
	return packEFace(type2.rtype, ptr)
}

func (type2 *unsafeIndirSliceType) UnsafeAppend(obj unsafe.Pointer, elem unsafe.Pointer) unsafe.Pointer {
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
