package native

import (
	"reflect"
	"github.com/v2pro/plz/accessor"
)

type intAccessor struct {
	accessor.NoopAccessor
	typ reflect.Type
}

func (acc *intAccessor) Kind() reflect.Kind {
	return reflect.Int
}

func (acc *intAccessor) GoString() string {
	return acc.typ.Name()
}

func (acc *intAccessor) Int(obj interface{}) int {
	return *((*int)(extractPtrFromEmptyInterface(obj)))
}

func (acc *intAccessor) SetInt(obj interface{}, val int) {
	if reflect.TypeOf(obj).Kind() != reflect.Ptr {
		panic("can only SetInt on pointer")
	}
	*((*int)(extractPtrFromEmptyInterface(obj))) = val
}
