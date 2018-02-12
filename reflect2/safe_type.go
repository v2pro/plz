package reflect2

import (
	"reflect"
	"unsafe"
)

type safeType struct {
	reflect.Type
	cfg *frozenConfig
}

func (type2 *safeType) New() interface{} {
	return reflect.New(type2.Type).Interface()
}

func (type2 *safeType) UnsafeNew() unsafe.Pointer {
	panic("does not support unsafe operation")
}

func (type2 *safeType) Elem() Type {
	return type2.cfg.Type2(type2.Type.Elem())
}

func (type2 *safeType) Type1() reflect.Type {
	return type2.Type
}

func (type2 *safeType) PackEFace(ptr unsafe.Pointer) interface{} {
	panic("does not support unsafe operation")
}