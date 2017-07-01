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
	case reflect.Interface:
		return &emptyInterfaceAccessor{}
	case reflect.Struct:
		return accessorOfStruct(typ)
	case reflect.Slice:
		templateElemObj := castToEmptyInterface(reflect.New(typ.Elem()).Interface())
		return &sliceAccessor{
			typ:              typ,
			templateElemObj:  templateElemObj,
		}
	case reflect.Array:
		templateElemObj := castToEmptyInterface(reflect.New(typ.Elem()).Interface())
		return &arrayAccessor{
			typ:              typ,
			templateElemObj:  templateElemObj,
		}
	}
	panic(fmt.Sprintf("do not support: %v", typ.Kind()))
}
