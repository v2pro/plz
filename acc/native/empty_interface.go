package native

import (
	"unsafe"
	"github.com/v2pro/plz/acc"
	"reflect"
	"github.com/v2pro/plz"
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

type ptrEmptyInterfaceAccessor struct {
	acc.NoopAccessor
}

func (accessor *ptrEmptyInterfaceAccessor) Kind() acc.Kind {
	return acc.Interface
}

func (accessor *ptrEmptyInterfaceAccessor) GoString() string {
	return "*interface{}"
}

func (accessor *ptrEmptyInterfaceAccessor) SetInt(obj interface{}, val int) {
	*(obj.(*interface{})) = val
}

func (accessor *ptrEmptyInterfaceAccessor) SetString(obj interface{}, val string) {
	*(obj.(*interface{})) = val
}

func (accessor *ptrEmptyInterfaceAccessor) Int(obj interface{}) int {
	obj = *(obj.(*interface{}))
	return plz.AccessorOf(reflect.TypeOf(obj)).Int(obj)
}

func (accessor *ptrEmptyInterfaceAccessor) String(obj interface{}) string {
	obj = *(obj.(*interface{}))
	return plz.AccessorOf(reflect.TypeOf(obj)).String(obj)
}

func (accessor *ptrEmptyInterfaceAccessor) PtrElem(obj interface{}) (interface{}, acc.Accessor) {
	obj = *(obj.(*interface{}))
	if obj == nil {
		return nil, nil
	}
	typ := reflect.TypeOf(obj)
	return obj, acc.AccessorOf(typ)
}

func (accessor *ptrEmptyInterfaceAccessor) SetPtrElem(obj interface{}, template interface{}) (elem interface{}, elemAccessor acc.Accessor) {
	typ := reflect.TypeOf(template)
	newObj := reflect.New(typ).Elem().Interface()
	*(obj.(*interface{})) = newObj
	return newObj, acc.AccessorOf(typ)
}
