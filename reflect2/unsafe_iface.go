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

func UnsafeIFaceToEFace(ptr unsafe.Pointer) interface{} {
	iface := (*iface)(ptr)
	if iface.itab == nil {
		return nil
	}
	return packEFace(iface.itab.rtype, iface.data)
}
