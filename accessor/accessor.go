package accessor

import (
	"reflect"
	"fmt"
)

var Providers = []func(reflect.Type) Accessor{}

func Of(typ reflect.Type) Accessor {
	for _, provider := range Providers {
		asor := provider(typ)
		if asor != nil {
			return asor
		}
	}
	panic(fmt.Sprintf("no accessor provider for: %v", typ))
}

type Accessor interface {
	Kind() reflect.Kind
	Int(obj interface{}) int
	SetInt(obj interface{}, val int)
	NumField() int
	Field(index int) StructField
	IterateMap(obj interface{}, cb func(key interface{}, value interface{}) bool)
	SetMapIndex(obj interface{}, key interface{}, value interface{})
	Key() Accessor
	Elem() Accessor
}

type StructField struct {
	Name     string
	Accessor Accessor
}

type NoopAccessor struct {
}

func (acc *NoopAccessor) Int(obj interface{}) int {
	panic("unsupported operation")
}

func (acc *NoopAccessor) SetInt(obj interface{}, val int) {
	panic("unsupported operation")
}

func (acc *NoopAccessor) NumField() int {
	panic("unsupported operation")
}

func (acc *NoopAccessor) Field(index int) StructField {
	panic("unsupported operation")
}

func (acc *NoopAccessor) IterateMap(obj interface{}, cb func(key interface{}, value interface{}) bool) {
	panic("unsupported operation")
}

func (acc *NoopAccessor) SetMapIndex(obj interface{}, key interface{}, value interface{}) {
	panic("unsupported operation")
}

func (acc *NoopAccessor) Key() Accessor {
	panic("unsupported operation")
}

func (acc *NoopAccessor) Elem() Accessor {
	panic("unsupported operation")
}
