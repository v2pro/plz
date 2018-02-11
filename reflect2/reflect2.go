package reflect2

import (
	"reflect"
	"unsafe"
)

//go:linkname unsafe_New reflect.unsafe_New
func unsafe_New(rtype unsafe.Pointer) unsafe.Pointer

type Type interface {
	reflect.Type
	// New return pointer to data of this type
	New() interface{}
	// UnsafeNew return the allocated space pointed by unsafe.Pointer
	UnsafeNew() unsafe.Pointer
}

type type2 struct {
	reflect.Type
	rtype  unsafe.Pointer
	prtype unsafe.Pointer
}

func (type2 type2) UnsafeNew() unsafe.Pointer {
	return unsafe_New(type2.rtype)
}

func (type2 type2) New() interface{} {
	return packEface(type2.prtype, type2.UnsafeNew())
}

func TypeOf(obj interface{}) Type {
	valType := reflect.TypeOf(obj)
	rtype := toEface(valType).data
	prtype := toEface(reflect.PtrTo(valType)).data
	return &type2{
		Type:   valType,
		rtype:  rtype,
		prtype: prtype,
	}
}

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
