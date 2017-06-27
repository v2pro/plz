package native

import (
	"github.com/v2pro/plz"
	"github.com/v2pro/plz/accessor"
	"reflect"
	"unsafe"
)

type sliceAccessor struct {
	accessor.NoopAccessor
	typ              reflect.Type
	templateSliceObj emptyInterface
	templateElemObj  emptyInterface
}

func (acc *sliceAccessor) Kind() reflect.Kind {
	return reflect.Slice
}

func (acc *sliceAccessor) Elem() accessor.Accessor {
	return plz.AccessorOf(reflect.PtrTo(acc.typ.Elem()))
}

func (acc *sliceAccessor) IterateArray(obj interface{}, cb func(elem interface{}) bool) {
	sliceHeader := extractSliceHeaderFromEmptyInterface(obj)
	elemSize := acc.typ.Elem().Size()
	for i := 0; i < sliceHeader.Len; i++ {
		elemPtr := uintptr(sliceHeader.Data) + uintptr(i)*elemSize
		elemObj := acc.templateElemObj
		elemObj.word = unsafe.Pointer(elemPtr)
		if !cb(castBackEmptyInterface(elemObj)) {
			return
		}
	}
}

func (acc *sliceAccessor) GrowOne(obj interface{}, elem interface{}) (interface{}, interface{}) {
	sliceHeader := extractSliceHeaderFromEmptyInterface(obj)
	at := sliceHeader.Len
	elemType := acc.typ.Elem()
	sliceHeader = growOne(sliceHeader, acc.typ, elemType)
	sliceObj := acc.templateSliceObj
	sliceObj.word = unsafe.Pointer(sliceHeader)
	elemPtr := uintptr(sliceHeader.Data) + uintptr(at)*elemType.Size()
	elemObj := acc.templateElemObj
	elemObj.word = unsafe.Pointer(elemPtr)
	return castBackEmptyInterface(sliceObj), castBackEmptyInterface(elemObj)
}

// grow grows the slice s so that it can hold extra more values, allocating
// more capacity if needed. It also returns the old and new slice lengths.
func growOne(slice *sliceHeader, sliceType reflect.Type, elementType reflect.Type) *sliceHeader {
	newLen := slice.Len + 1
	if newLen <= slice.Cap {
		slice.Len = newLen
		return slice
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
	return &sliceHeader{
		Len:  newLen,
		Cap:  newCap,
		Data: dst,
	}
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
