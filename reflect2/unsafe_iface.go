package reflect2

import "unsafe"

type iface struct {
	itab *itab
	data unsafe.Pointer
}

type itab struct {
	ignore unsafe.Pointer
	rtype  unsafe.Pointer
}
