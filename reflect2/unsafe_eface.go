package reflect2

import "unsafe"

type eface struct {
	rtype unsafe.Pointer
	data  unsafe.Pointer
}

func toEface(obj interface{}) *eface {
	return (*eface)(unsafe.Pointer(&obj))
}

func packEface(rtype unsafe.Pointer, data unsafe.Pointer) interface{} {
	var i interface{}
	e := (*eface)(unsafe.Pointer(&i))
	e.rtype = rtype
	e.data = data
	return i
}