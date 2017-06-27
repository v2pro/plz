package native

import (
	"github.com/v2pro/plz/accessor"
	"reflect"
	"fmt"
)

func init() {
	accessor.Providers = append(accessor.Providers, accessorOf)
}

func accessorOf(typ reflect.Type) accessor.Accessor {
	if typ.Kind() == reflect.Map {
		return &mapAccessor{
			typ: typ,
		}
	}
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	switch typ.Kind() {
	case reflect.Int:
		return &intAccessor{}
	case reflect.Struct:
		return &structAccessor{
			typ: typ,
		}
	case reflect.Slice:
		return &sliceAccessor{
			typ: typ,
			templateSliceObj: castToEmptyInterface(reflect.New(typ).Elem().Interface()),
			templateElemObj: castToEmptyInterface(reflect.New(typ.Elem()).Interface()),
		}
	}
	panic(fmt.Sprintf("do not support: %v", typ.Kind()))
}

type intAccessor struct {
	accessor.NoopAccessor
}

func (acc *intAccessor) Kind() reflect.Kind {
	return reflect.Int
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
