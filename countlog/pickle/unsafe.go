package pickle

import (
	"unsafe"
)

// emptyInterface is the header for an interface{} value.
type emptyInterface struct {
	typ  unsafe.Pointer
	word unsafe.Pointer
}

type stringWritableHeader struct {
	Data uintptr
	Len  int
}

type sliceReadonlyHeader struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}

type sliceWritableHeader struct {
	Data uintptr
	Len  int
	Cap  int
}

func ptrOfEmptyInterface(obj interface{}) unsafe.Pointer {
	return unsafe.Pointer((*emptyInterface)(unsafe.Pointer(&obj)).word)
}

func ptrAsBytes(size int, ptr unsafe.Pointer) []byte {
	valAsSlice := *(*[]byte)((unsafe.Pointer)(&sliceReadonlyHeader{
		Data: ptr, Len: size, Cap: size}))
	return valAsSlice
}
