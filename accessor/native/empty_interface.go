package native

import "unsafe"

func castToEmptyInterface(val interface{}) emptyInterface {
	return *((*emptyInterface)(unsafe.Pointer(&val)))
}

func castBackEmptyInterface(ei emptyInterface) interface{} {
	return *((*interface{})(unsafe.Pointer(&ei)))
}

func extractPtrFromEmptyInterface(val interface{}) unsafe.Pointer {
	return castToEmptyInterface(val).word
}

// emptyInterface is the header for an interface{} value.
type emptyInterface struct {
	typ  unsafe.Pointer
	word unsafe.Pointer
}
