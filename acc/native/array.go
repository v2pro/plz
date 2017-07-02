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

func (accessor *arrayAccessor) IterateArray(obj interface{}, cb func(index int, elem interface{}) bool) {
	elemSize := accessor.typ.Elem().Size()
	head := uintptr(extractPtrFromEmptyInterface(obj))
	for index := 0; index < accessor.typ.Len(); index++ {
		elemPtr := head + uintptr(index) * elemSize
		elemObj := accessor.templateElemObj
		elemObj.word = unsafe.Pointer(elemPtr)
		if !cb(index, castBackEmptyInterface(elemObj)) {
			return
		}
	}
}

func (accessor *arrayAccessor) FillArray(obj interface{}, cb func(filler acc.ArrayFiller)) {
	elemSize := accessor.typ.Elem().Size()
	head := uintptr(extractPtrFromEmptyInterface(obj))
	filler := &arrayFiller{
		elemSize:        elemSize,
		len:             accessor.typ.Len(),
		head:            head,
		templateElemObj: accessor.templateElemObj,
	}
	cb(filler)
}

type arrayFiller struct {
	index           int
	len             int
	elemSize        uintptr
	head            uintptr
	templateElemObj emptyInterface
}

func (filler *arrayFiller) Next() (int, interface{}) {
	if filler.index < filler.len {
		elemObj := filler.templateElemObj
		elemObj.word = unsafe.Pointer(filler.head + uintptr(filler.index)*filler.elemSize)
		currentIndex := filler.index
		filler.index++
		return currentIndex, castBackEmptyInterface(elemObj)
	} else {
		return -1, nil
	}
}

func (filler *arrayFiller) Fill() {
}
