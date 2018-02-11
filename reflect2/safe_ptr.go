package reflect2

import (
	"reflect"
	"unsafe"
)

type safePtrType struct {
	safeType
}

func (type2 *safePtrType) Get(obj interface{}) interface{} {
	return reflect.ValueOf(obj).Elem().Interface()
}

func (type2 *safePtrType) UnsafeGet(obj unsafe.Pointer) unsafe.Pointer {
	panic("does not support unsafe operation")
}
