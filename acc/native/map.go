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
		value := reflectVal.MapIndex(key).Interface()
		if accessor.typ.Elem().Kind() == reflect.Interface {
			if !cb(key.Interface(), &value) {
				return
			}
		} else {
			if !cb(key.Interface(), value) {
				return
			}
		}
	}
}

func (accessor *mapAccessor) FillMap(obj interface{}, cb func(filler acc.MapFiller)) {
	filler := &mapFiller{
		typ: accessor.typ,
		value: reflect.ValueOf(obj),
	}
	cb(filler)
}

type mapFiller struct {
	typ   reflect.Type
	value reflect.Value
	lastKey reflect.Value
	lastElem reflect.Value
}

func (filler *mapFiller) Next() (interface{}, interface{}) {
	filler.lastKey = reflect.New(filler.typ.Key())
	filler.lastElem = reflect.New(filler.typ.Elem())
	return filler.lastKey.Interface(), filler.lastElem.Interface()
}

func (filler *mapFiller) Fill() {
	filler.value.SetMapIndex(filler.lastKey.Elem(), filler.lastElem.Elem())
}

func (accessor *mapAccessor) Key() acc.Accessor {
	return plz.AccessorOf(accessor.typ.Key())
}

func (accessor *mapAccessor) Elem() acc.Accessor {
	return plz.AccessorOf(accessor.typ.Elem())
}
