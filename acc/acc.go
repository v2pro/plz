package acc

import (
	"fmt"
	"reflect"
)

var Providers = []func(typ reflect.Type) Accessor{}

func AccessorOf(typ reflect.Type) Accessor {
	for _, provider := range Providers {
		asor := provider(typ)
		if asor != nil {
			return asor
		}
	}
	panic(fmt.Sprintf("no accessor provider for: %v", typ))
}

type Kind uint

const (
	Invalid Kind = iota
	Bool
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Uintptr
	Float32
	Float64
	Array
	Interface
	Map
	String
	Struct
)

type Accessor interface {
	fmt.GoStringer
	Kind() Kind
	KindOf(obj interface{}) Kind
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
	AppendArray(obj interface{}, setElem func(elem interface{}))
	// primitives
	Skip(obj interface{})
	Int(obj interface{}) int
	SetInt(obj interface{}, val int)
	String(obj interface{}) string
	SetString(obj interface{}, val string)
	Uintptr(obj interface{}) uintptr
	Interface(obj interface{}) interface{}
}

type StructField struct {
	Name     string
	Accessor Accessor
	Tags     map[string]interface{}
}

type NoopAccessor struct {
}

func (acc *NoopAccessor) KindOf(obj interface{}) Kind {
	panic("unsupported operation")
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

func (acc *NoopAccessor) AppendArray(obj interface{}, setElem func(elem interface{})) {
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

func (acc *NoopAccessor) Uintptr(obj interface{}) uintptr {
	panic("unsupported operation")
}

func (acc *NoopAccessor) Skip(obj interface{}) {
}

func (acc *NoopAccessor) Interface(obj interface{}) interface{} {
	panic("unsupported operation")
}
