package native

import (
	"reflect"
	"github.com/v2pro/plz/acc"
)

type intAccessor struct {
	acc.NoopAccessor
	typ reflect.Type
}

func (accessor *intAccessor) Kind() acc.Kind {
	return acc.Int
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
