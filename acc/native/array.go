package native

import (
	"github.com/v2pro/plz/acc"
	"reflect"
	"github.com/v2pro/plz"
	"unsafe"
)

type arrayAccessor struct {
	acc.NoopAccessor
	typ             reflect.Type
	templateElemObj emptyInterface
}

func (accessor *arrayAccessor) Kind() acc.Kind {
	return acc.Array
}

func (accessor *arrayAccessor) GoString() string {
	return accessor.typ.String()
}

func (accessor *arrayAccessor) Elem() acc.Accessor {
	return plz.AccessorOf(reflect.PtrTo(accessor.typ.Elem()))
}

func (accessor *arrayAccessor) IterateArray(obj interface{}, cb func(elem interface{}) bool) {
	ptr := uintptr(extractPtrFromEmptyInterface(obj))
	elemSize := accessor.typ.Elem().Size()
	for i := 0; i < accessor.typ.Len(); i++ {
		elemPtr := ptr + uintptr(i)*elemSize
		elemObj := accessor.templateElemObj
		elemObj.word = unsafe.Pointer(elemPtr)
		if !cb(castBackEmptyInterface(elemObj)) {
			return
		}
	}
}
