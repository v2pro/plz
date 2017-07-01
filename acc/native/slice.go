package native

import (
	"github.com/v2pro/plz"
	"github.com/v2pro/plz/acc"
	"reflect"
	"unsafe"
)

type sliceAccessor struct {
	acc.NoopAccessor
	typ              reflect.Type
	templateElemObj  emptyInterface
}

func (accessor *sliceAccessor) Kind() acc.Kind {
	return acc.Array
}

func (accessor *sliceAccessor) GoString() string {
	return accessor.typ.String()
}

func (accessor *sliceAccessor) Elem() acc.Accessor {
	return plz.AccessorOf(reflect.PtrTo(accessor.typ.Elem()))
}

func (accessor *sliceAccessor) IterateArray(obj interface{}, cb func(elem interface{}) bool) {
	sliceHeader := extractSliceHeaderFromEmptyInterface(obj)
	elemSize := accessor.typ.Elem().Size()
	for i := 0; i < sliceHeader.Len; i++ {
		elemPtr := uintptr(sliceHeader.Data) + uintptr(i)*elemSize
		elemObj := accessor.templateElemObj
		elemObj.word = unsafe.Pointer(elemPtr)
		if !cb(castBackEmptyInterface(elemObj)) {
			return
		}
	}
}

func (accessor *sliceAccessor) AppendArray(obj interface{}, setElem func(elem interface{})) {
	sliceHeader := extractSliceHeaderFromEmptyInterface(obj)
	at := sliceHeader.Len
	elemType := accessor.typ.Elem()
	growOne(sliceHeader, accessor.typ, elemType)
	elemPtr := uintptr(sliceHeader.Data) + uintptr(at)*elemType.Size()
	elemObj := accessor.templateElemObj
	elemObj.word = unsafe.Pointer(elemPtr)
	setElem(castBackEmptyInterface(elemObj))
}

// grow grows the slice s so that it can hold extra more values, allocating
// more capacity if needed. It also returns the old and new slice lengths.
func growOne(slice *sliceHeader, sliceType reflect.Type, elementType reflect.Type) {
	newLen := slice.Len + 1
	if newLen <= slice.Cap {
		slice.Len = newLen
		return
	}
	newCap := slice.Cap
	if newCap == 0 {
		newCap = 1
	} else {
		for newCap < newLen {
			if slice.Len < 1024 {
				newCap += newCap
			} else {
				newCap += newCap / 4
			}
		}
	}
	dst := unsafe.Pointer(reflect.MakeSlice(sliceType, newLen, newCap).Pointer())
	// copy old array into new array
	originalBytesCount := uintptr(slice.Len) * elementType.Size()
	srcPtr := (*[1 << 30]byte)(slice.Data)
	dstPtr := (*[1 << 30]byte)(dst)
	for i := uintptr(0); i < originalBytesCount; i++ {
		dstPtr[i] = srcPtr[i]
	}
	slice.Len = newLen
	slice.Cap = newCap
	slice.Data = dst
}

// sliceHeader is a safe version of SliceHeader used within this package.
type sliceHeader struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}

func extractSliceHeaderFromEmptyInterface(obj interface{}) *sliceHeader {
	ptr := extractPtrFromEmptyInterface(obj)
	return (*sliceHeader)(ptr)
}
