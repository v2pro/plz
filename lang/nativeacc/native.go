package nativeacc

import (
	"fmt"
	"github.com/v2pro/plz/lang"
	"reflect"
)

func init() {
	lang.Providers = append(lang.Providers, accessorOfNativeType)
}

func accessorOfNativeType(typ reflect.Type) lang.Accessor {
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
				NoopAccessor:  lang.NoopAccessor{"ptrIntAccessor"},
				valueAccessor: lang.AccessorOf(elemType),
			}}
		case reflect.String:
			return &ptrStringAccessor{ptrAccessor{
				NoopAccessor:  lang.NoopAccessor{"ptrStringAccessor"},
				valueAccessor: lang.AccessorOf(elemType),
			}}
		case reflect.Interface:
			return &ptrEmptyInterfaceAccessor{
				lang.NoopAccessor{"ptrEmptyInterfaceAccessor"}}
		case reflect.Struct:
			fallthrough
		case reflect.Slice:
			fallthrough
		case reflect.Array:
			return accessorOf(elemType)

		}
	case reflect.Int:
		return &intAccessor{
			NoopAccessor: lang.NoopAccessor{"intAccessor"},
			typ:          typ,
		}
	case reflect.String:
		return &stringAccessor{
			NoopAccessor: lang.NoopAccessor{"stringAccessor"},
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
