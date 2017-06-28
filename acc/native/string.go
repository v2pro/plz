package native

import (
	"reflect"
	"github.com/v2pro/plz/acc"
)

type stringAccessor struct {
	acc.NoopAccessor
	typ reflect.Type
}

func (accessor *stringAccessor) Kind() reflect.Kind {
	return reflect.String
}

func (accessor *stringAccessor) GoString() string {
	return accessor.typ.Name()
}

func (accessor *stringAccessor) String(obj interface{}) string {
	return *((*string)(extractPtrFromEmptyInterface(obj)))
}

func (accessor *stringAccessor) SetString(obj interface{}, val string) {
	if reflect.TypeOf(obj).Kind() != reflect.Ptr {
		panic("can only SetString on pointer")
	}
	*((*string)(extractPtrFromEmptyInterface(obj))) = val
}