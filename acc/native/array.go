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
	elemSize := accessor.typ.Elem().Size()
	head := uintptr(extractPtrFromEmptyInterface(obj))
	tail := head + uintptr(accessor.typ.Len())*elemSize
	for elemPtr := head; elemPtr < tail; elemPtr += elemSize {
		elemObj := accessor.templateElemObj
		elemObj.word = unsafe.Pointer(elemPtr)
		if !cb(castBackEmptyInterface(elemObj)) {
			return
		}
	}
}

func (accessor *arrayAccessor) FillArray(obj interface{}, cb func(filler acc.ArrayFiller)) {
	elemSize := accessor.typ.Elem().Size()
	head := uintptr(extractPtrFromEmptyInterface(obj))
	tail := head + uintptr(accessor.typ.Len())*elemSize
	elemPtr := head
	filler := &arrayFiller{
		elemSize:        elemSize,
		elemPtr:         elemPtr,
		tail:            tail,
		templateElemObj: accessor.templateElemObj,
	}
	cb(filler)
}

type arrayFiller struct {
	elemSize        uintptr
	elemPtr         uintptr
	tail            uintptr
	templateElemObj emptyInterface
}

func (filler *arrayFiller) Next() interface{} {
	if filler.elemPtr < filler.tail {
		elemObj := filler.templateElemObj
		elemObj.word = unsafe.Pointer(filler.elemPtr)
		filler.elemPtr += filler.elemSize
		return castBackEmptyInterface(elemObj)
	} else {
		return nil
	}
}

func (filler *arrayFiller) Fill() {
}
