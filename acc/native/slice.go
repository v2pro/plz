package native

import (
	"github.com/v2pro/plz"
	"github.com/v2pro/plz/acc"
	"reflect"
	"unsafe"
)

type sliceAccessor struct {
	acc.NoopAccessor
	typ             reflect.Type
	templateElemObj emptyInterface
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

func (accessor *sliceAccessor) IterateArray(obj interface{}, cb func(index int, elem interface{}) bool) {
	sliceHeader := extractSliceHeaderFromEmptyInterface(obj)
	elemSize := accessor.typ.Elem().Size()
	head := uintptr(sliceHeader.Data)
	for index := 0; index < sliceHeader.Len; index++ {
		elemPtr := head + uintptr(index)*elemSize
		elemObj := accessor.templateElemObj
		elemObj.word = unsafe.Pointer(elemPtr)
		if !cb(index, castBackEmptyInterface(elemObj)) {
			return
		}
	}
}

func (accessor *sliceAccessor) FillArray(obj interface{}, cb func(filler acc.ArrayFiller)) {
	sliceHeader := extractSliceHeaderFromEmptyInterface(obj)
	sliceHeader.Len = 0
	filler := &sliceFiller{
		sliceTyp:        accessor.typ,
		elemType:        accessor.typ.Elem(),
		sliceHeader:     sliceHeader,
		templateElemObj: accessor.templateElemObj,
	}
	cb(filler)
}

type sliceFiller struct {
	sliceTyp        reflect.Type
	elemType        reflect.Type
	sliceHeader     *sliceHeader
	templateElemObj emptyInterface
}

func (filler *sliceFiller) Next() (int, interface{}) {
	header := filler.sliceHeader
	at := header.Len
	growOne(header, filler.sliceTyp, filler.elemType)
	elemObj := filler.templateElemObj
	elemPtr := uintptr(header.Data) + uintptr(at)*filler.elemType.Size()
	elemObj.word = unsafe.Pointer(elemPtr)
	return at, castBackEmptyInterface(elemObj)
}

func (filler *sliceFiller) Fill() {
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
