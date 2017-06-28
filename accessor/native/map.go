package native

import (
	"github.com/v2pro/plz"
	"github.com/v2pro/plz/accessor"
	"reflect"
)

type mapAccessor struct {
	accessor.NoopAccessor
	typ reflect.Type
}

func (acc *mapAccessor) Kind() reflect.Kind {
	return reflect.Map
}

func (acc *mapAccessor) IterateMap(obj interface{}, cb func(key interface{}, value interface{}) bool) {
	reflectVal := reflect.ValueOf(obj)
	for _, key := range reflectVal.MapKeys() {
		value := reflectVal.MapIndex(key)
		if !cb(key.Interface(), value.Interface()) {
			return
		}
	}
}

func (acc *mapAccessor) SetMap(obj interface{}, setKey func(key interface{}), setElem func(elem interface{})) {
	key := reflect.New(acc.typ.Key())
	setKey(key.Interface())
	elem := reflect.New(acc.typ.Elem())
	setElem(elem.Interface())
	reflect.ValueOf(obj).SetMapIndex(key.Elem(), elem.Elem())
}

func (acc *mapAccessor) Key() accessor.Accessor {
	return plz.AccessorOf(acc.typ.Key())
}

func (acc *mapAccessor) Elem() accessor.Accessor {
	return plz.AccessorOf(acc.typ.Elem())
}
