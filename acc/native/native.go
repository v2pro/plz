package native

import (
	"fmt"
	"github.com/v2pro/plz/acc"
	"reflect"
	"github.com/v2pro/plz/tagging"
)

func init() {
	acc.Providers = append(acc.Providers, accessorOf)
}

func accessorOf(typ reflect.Type, profile string) acc.Accessor {
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
		tags := tagging.Get(typ)
		return &structAccessor{
			typ: typ,
			tags: tags,
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
