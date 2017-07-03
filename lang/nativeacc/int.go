package nativeacc

import (
	"reflect"
	"github.com/v2pro/plz/lang"
)

type intAccessor struct {
	lang.NoopAccessor
	typ reflect.Type
}

func (accessor *intAccessor) Kind() lang.Kind {
	return lang.Int
}

func (accessor *intAccessor) GoString() string {
	return accessor.typ.Name()
}

func (accessor *intAccessor) Int(obj interface{}) int {
	return *((*int)(extractPtrFromEmptyInterface(obj)))
}

type ptrIntAccessor struct {
	ptrAccessor
}

func (accessor *ptrIntAccessor) Int(obj interface{}) int {
	return accessor.valueAccessor.Int(obj)
}

func (accessor *ptrIntAccessor) SetInt(obj interface{}, val int) {
	*((*int)(extractPtrFromEmptyInterface(obj))) = val
}
