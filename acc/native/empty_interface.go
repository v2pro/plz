package native

import (
	"unsafe"
	"github.com/v2pro/plz/acc"
	"reflect"
	"fmt"
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

func (accessor *emptyInterfaceAccessor) KindOf(obj interface{}) acc.Kind {
	obj = *(obj.(*interface{}))
	typ := reflect.TypeOf(obj)
	switch typ.Kind() {
	case reflect.Int:
		return acc.Int
	case reflect.String:
		return acc.String
	}
	panic(fmt.Sprintf("KindOf does not support: %v", typ))
}

func (accessor *emptyInterfaceAccessor) Kind() acc.Kind {
	return acc.Interface
}

func (accessor *emptyInterfaceAccessor) GoString() string {
	return "interface{}"
}

func (accessor *emptyInterfaceAccessor) Int(obj interface{}) int {
	obj = *(obj.(*interface{}))
	return *((*int)(extractPtrFromEmptyInterface(obj)))
}

func (accessor *emptyInterfaceAccessor) SetInt(obj interface{}, val int) {
	objPtr := obj.(*interface{})
	*objPtr = val
}

func (accessor *emptyInterfaceAccessor) String(obj interface{}) string {
	obj = *(obj.(*interface{}))
	return *((*string)(extractPtrFromEmptyInterface(obj)))
}
