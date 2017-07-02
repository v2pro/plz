package native

import (
	"reflect"
	"github.com/v2pro/plz/acc"
)

type stringAccessor struct {
	acc.NoopAccessor
	typ reflect.Type
}

func (accessor *stringAccessor) Kind() acc.Kind {
	return acc.String
}

func (accessor *stringAccessor) GoString() string {
	return accessor.typ.Name()
}

func (accessor *stringAccessor) String(obj interface{}) string {
	return *((*string)(extractPtrFromEmptyInterface(obj)))
}


type ptrStringAccessor struct {
	ptrAccessor
}

func (accessor *ptrStringAccessor) String(obj interface{}) string {
	return accessor.valueAccessor.String(obj)
}

func (accessor *ptrStringAccessor) SetString(obj interface{}, val string) {
	*((*string)(extractPtrFromEmptyInterface(obj))) = val
}
