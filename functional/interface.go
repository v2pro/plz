package functional

import "unsafe"

// emptyInterface is the header for an interface{} value.
type emptyInterface struct {
	typ  unsafe.Pointer
	word unsafe.Pointer
}

func toPointer(val interface{}) unsafe.Pointer {
	return (*((*emptyInterface)(unsafe.Pointer(&val)))).word
}
