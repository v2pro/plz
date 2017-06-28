package accessor

import (
	"fmt"
	"reflect"
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
	// map
	Key() Accessor
	// array/map
	Elem() Accessor
	// struct
	NumField() int
	Field(index int) StructField
	// map
	IterateMap(obj interface{}, cb func(key interface{}, elem interface{}) bool)
	SetMap(obj interface{}, setKey func(key interface{}), setElem func(elem interface{}))
	// array
	IterateArray(obj interface{}, cb func(elem interface{}) bool)
	AppendArray(obj interface{}, setElem func(elem interface{})) interface{}
	// primitives
	Int(obj interface{}) int
	SetInt(obj interface{}, val int)
	String(obj interface{}) string
	SetString(obj interface{}, val string)
}

type StructField struct {
	Name     string
	Accessor Accessor
}

type NoopAccessor struct {
}

func (acc *NoopAccessor) NumField() int {
	panic("unsupported operation")
}

func (acc *NoopAccessor) Field(index int) StructField {
	panic("unsupported operation")
}

func (acc *NoopAccessor) IterateMap(obj interface{}, cb func(key interface{}, elem interface{}) bool) {
	panic("unsupported operation")
}

func (acc *NoopAccessor) SetMap(obj interface{}, setKey func(key interface{}), setElem func(elem interface{})) {
	panic("unsupported operation")
}

func (acc *NoopAccessor) Key() Accessor {
	panic("unsupported operation")
}

func (acc *NoopAccessor) Elem() Accessor {
	panic("unsupported operation")
}

func (acc *NoopAccessor) IterateArray(obj interface{}, cb func(elem interface{}) bool) {
	panic("unsupported operation")
}

func (acc *NoopAccessor) AppendArray(obj interface{}, setElem func(elem interface{})) interface{} {
	panic("unsupported operation")
}

func (acc *NoopAccessor) Int(obj interface{}) int {
	panic("unsupported operation")
}

func (acc *NoopAccessor) SetInt(obj interface{}, val int) {
	panic("unsupported operation")
}

func (acc *NoopAccessor) String(obj interface{}) string {
	panic("unsupported operation")
}

func (acc *NoopAccessor) SetString(obj interface{}, val string) {
	panic("unsupported operation")
}
