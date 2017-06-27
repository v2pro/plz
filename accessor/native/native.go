package native

import (
	"fmt"
	"github.com/v2pro/plz/accessor"
	"reflect"
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
	case reflect.String:
		return &stringAccessor{}
	case reflect.Struct:
		return &structAccessor{
			typ: typ,
		}
	case reflect.Slice:
		return &sliceAccessor{
			typ:              typ,
			templateSliceObj: castToEmptyInterface(reflect.New(typ).Elem().Interface()),
			templateElemObj:  castToEmptyInterface(reflect.New(typ.Elem()).Interface()),
		}
	}
	panic(fmt.Sprintf("do not support: %v", typ.Kind()))
}
