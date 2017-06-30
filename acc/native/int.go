package native

import (
	"reflect"
	"github.com/v2pro/plz/acc"
)

type intAccessor struct {
	acc.NoopAccessor
	typ reflect.Type
}

func (accessor *intAccessor) Kind() reflect.Kind {
	return reflect.Int
}

func (accessor *intAccessor) GoString() string {
	return accessor.typ.Name()
}

func (accessor *intAccessor) Int(obj interface{}) int {
	return *((*int)(extractPtrFromEmptyInterface(obj)))
}

func (accessor *intAccessor) SetInt(obj interface{}, val int) {
	if reflect.TypeOf(obj).Kind() != reflect.Ptr {
		panic("can only SetInt on pointer")
	}
	*((*int)(extractPtrFromEmptyInterface(obj))) = val
}
