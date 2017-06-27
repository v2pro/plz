package native

import (
	"github.com/v2pro/plz/accessor"
	"reflect"
	"fmt"
)

func init() {
	accessor.Providers = append(accessor.Providers, accessorOf)
}

func accessorOf(obj interface{}) accessor.Accessor {
	typ := reflect.TypeOf(obj)
	if typ.Kind() != reflect.Ptr {
		return nil
	}
	typ = typ.Elem()
	switch typ.Kind() {
	case reflect.Int:
		return &intAccessor{}
	}
	panic(fmt.Sprintf("do not support: %v", typ.Kind()))
}

type intAccessor struct {
}

func (accessor *intAccessor) Kind() reflect.Kind {
	return reflect.Int
}

func (accessor *intAccessor) Int(obj interface{}) int {
	typedObj := obj.(*int)
	return *typedObj
}

func (accessor *intAccessor) SetInt(obj interface{}, val int) {
	*(obj.(*int)) = val
}
