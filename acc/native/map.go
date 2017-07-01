package native

import (
	"github.com/v2pro/plz"
	"github.com/v2pro/plz/acc"
	"reflect"
)

type mapAccessor struct {
	acc.NoopAccessor
	typ reflect.Type
}

func (accessor *mapAccessor) Kind() acc.Kind {
	return acc.Map
}

func (accessor *mapAccessor) GoString() string {
	return accessor.typ.String()
}

func (accessor *mapAccessor) IterateMap(obj interface{}, cb func(key interface{}, value interface{}) bool) {
	reflectVal := reflect.ValueOf(obj)
	for _, key := range reflectVal.MapKeys() {
		value := reflectVal.MapIndex(key)
		if !cb(key.Interface(), value.Interface()) {
			return
		}
	}
}

func (accessor *mapAccessor) SetMap(obj interface{}, cb func(key interface{}, elem interface{})) {
	key := reflect.New(accessor.typ.Key())
	elem := reflect.New(accessor.typ.Elem())
	cb(key.Interface(), elem.Interface())
	mapVal := reflect.ValueOf(obj)
	mapVal.SetMapIndex(key.Elem(), elem.Elem())
}

func (accessor *mapAccessor) Key() acc.Accessor {
	return plz.AccessorOf(accessor.typ.Key())
}

func (accessor *mapAccessor) Elem() acc.Accessor {
	return plz.AccessorOf(accessor.typ.Elem())
}
