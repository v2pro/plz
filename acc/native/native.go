package native

import (
	"fmt"
	"github.com/v2pro/plz/acc"
	"reflect"
)

func init() {
	acc.Providers = append(acc.Providers, accessorOf)
}

func accessorOf(typ reflect.Type) acc.Accessor {
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
		return &intAccessor{typ:typ}
	case reflect.String:
		return &stringAccessor{typ:typ}
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
