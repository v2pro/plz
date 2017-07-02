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
	switch typ.Kind() {
	case reflect.Ptr:
		elemType := typ.Elem()
		switch elemType.Kind() {
		case reflect.Int:
			return &ptrIntAccessor{ptrAccessor{
				NoopAccessor:  acc.NoopAccessor{"ptrIntAccessor"},
				valueAccessor: acc.AccessorOf(elemType),
			}}
		case reflect.String:
			return &ptrStringAccessor{ptrAccessor{
				NoopAccessor:  acc.NoopAccessor{"ptrStringAccessor"},
				valueAccessor: acc.AccessorOf(elemType),
			}}
		case reflect.Interface:
			return &ptrEmptyInterfaceAccessor{
				acc.NoopAccessor{"ptrEmptyInterfaceAccessor"}}
		case reflect.Struct:
			fallthrough
		case reflect.Slice:
			fallthrough
		case reflect.Array:
			return accessorOf(elemType)

		}
	case reflect.Int:
		return &intAccessor{
			NoopAccessor: acc.NoopAccessor{"intAccessor"},
			typ: typ,
		}
	case reflect.String:
		return &stringAccessor{
			NoopAccessor: acc.NoopAccessor{"stringAccessor"},
			typ:          typ,
		}
	case reflect.Struct:
		return accessorOfStruct(typ)
	case reflect.Slice:
		templateElemObj := castToEmptyInterface(reflect.New(typ.Elem()).Interface())
		return &sliceAccessor{
			typ:             typ,
			templateElemObj: templateElemObj,
		}
	case reflect.Array:
		templateElemObj := castToEmptyInterface(reflect.New(typ.Elem()).Interface())
		return &arrayAccessor{
			typ:             typ,
			templateElemObj: templateElemObj,
		}
	}
	panic(fmt.Sprintf("do not support: %v", typ.Kind()))
}
