package native

import (
	"unsafe"
	"github.com/v2pro/plz/acc"
	"reflect"
)

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

type emptyInterfaceAccessor struct {
	acc.NoopAccessor
}

func (accessor *emptyInterfaceAccessor) AccessorOf(obj interface{}) acc.Accessor {
	obj = *(obj.(*interface{}))
	typ := reflect.TypeOf(obj)
	return &deferenceAccessor{realAcc:acc.AccessorOf(typ)}
}

func (accessor *emptyInterfaceAccessor) Kind() acc.Kind {
	return acc.Interface
}

func (accessor *emptyInterfaceAccessor) Key() acc.Accessor {
	return &stringAccessor{}
}

func (accessor *emptyInterfaceAccessor) Elem() acc.Accessor {
	return accessor
}

func (accessor *emptyInterfaceAccessor) GoString() string {
	return "interface{}"
}

func (accessor *emptyInterfaceAccessor) SetInt(obj interface{}, val int) {
	objPtr := obj.(*interface{})
	*objPtr = val
}

func (accessor *emptyInterfaceAccessor) SetString(obj interface{}, val string) {
	objPtr := obj.(*interface{})
	*objPtr = val
}

func (accessor *emptyInterfaceAccessor) FillMap(obj interface{}, cb func(filler acc.MapFiller)) {
	realObj := obj.(*interface{})
	if *realObj == nil {
		*realObj = map[string]interface{}{}
	}
	m := (*realObj).(map[string]interface{})
	filler := &genericMapFiller{
		m: m,
	}
	cb(filler)
}

type genericMapFiller struct {
	m        map[string]interface{}
	lastKey  string
	lastElem interface{}
}

func (filler *genericMapFiller) Next() (interface{}, interface{}) {
	return &filler.lastKey, &filler.lastElem
}

func (filler *genericMapFiller) Fill() {
	filler.m[filler.lastKey] = filler.lastElem
}

func (accessor *emptyInterfaceAccessor) FillArray(obj interface{}, cb func(filler acc.ArrayFiller)) {
	realObj := obj.(*interface{})
	if *realObj == nil {
		*realObj = []interface{}{}
	}
	arr := (*realObj).([]interface{})
	filler := &genericArrayFiller{
		arr: arr,
	}
	cb(filler)
	*realObj = filler.arr
}

type genericArrayFiller struct {
	arr      []interface{}
	lastElem interface{}
}

func (filler *genericArrayFiller) Next() interface{} {
	return &filler.lastElem
}

func (filler *genericArrayFiller) Fill() {
	filler.arr = append(filler.arr, filler.lastElem)
}
