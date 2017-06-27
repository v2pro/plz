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

func (acc *mapAccessor) SetMapIndex(obj interface{}, key interface{}, value interface{}) {
	reflect.ValueOf(obj).SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(value))
}

func (acc *mapAccessor) Key() accessor.Accessor {
	return plz.AccessorOf(acc.typ.Key())
}

func (acc *mapAccessor) Elem() accessor.Accessor {
	return plz.AccessorOf(acc.typ.Elem())
}
