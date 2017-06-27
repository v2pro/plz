package native

import (
	"github.com/v2pro/plz/accessor"
	"reflect"
)

type stringAccessor struct {
	accessor.NoopAccessor
}

func (acc *stringAccessor) Kind() reflect.Kind {
	return reflect.Int
}

func (acc *stringAccessor) String(obj interface{}) string {
	return *((*string)(extractPtrFromEmptyInterface(obj)))
}

func (acc *stringAccessor) SetString(obj interface{}, val string) {
	if reflect.TypeOf(obj).Kind() != reflect.Ptr {
		panic("can only SetString on pointer")
	}
	*((*string)(extractPtrFromEmptyInterface(obj))) = val
}