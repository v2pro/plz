package acc

import (
	"fmt"
	"reflect"
)

var Providers = []func(typ reflect.Type) Accessor{}

func AccessorOf(typ reflect.Type) Accessor {
	for _, provider := range Providers {
		accessor := provider(typ)
		if accessor != nil {
			return accessor
		}
	}
	panic(fmt.Sprintf("no accessor provider for: %v", typ))
}

type Kind uint

const (
	Invalid   Kind = iota
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

func (kind Kind) IsSingleValue() bool {
	switch kind {
	case Invalid:
		return true
	case Bool:
		return true
	case Int:
		return true
	case Int8:
		return true
	case Int16:
		return true
	case Int32:
		return true
	case Int64:
		return true
	case Uint:
		return true
	case Uint8:
		return true
	case Uint16:
		return true
	case Uint32:
		return true
	case Uint64:
		return true
	case Uintptr:
		return true
	case Float32:
		return true
	case Float64:
		return true
	case Array:
		return false
	case Interface:
		return false
	case Map:
		return false
	case String:
		return true
	case Struct:
		return false
	}
	return false
}

func (kind Kind) GoString() string {
	switch kind {
	case Invalid:
		return "Invalid"
	case Bool:
		return "Bool"
	case Int:
		return "Int"
	case Int8:
		return "Int8"
	case Int16:
		return "Int16"
	case Int32:
		return "Int32"
	case Int64:
		return "Int64"
	case Uint:
		return "Uint"
	case Uint8:
		return "Uint8"
	case Uint16:
		return "Uint16"
	case Uint32:
		return "Uint32"
	case Uint64:
		return "Uint64"
	case Uintptr:
		return "Uintptr"
	case Float32:
		return "Float32"
	case Float64:
		return "Float64"
	case Array:
		return "Array"
	case Interface:
		return "Interface"
	case Map:
		return "Map"
	case String:
		return "String"
	case Struct:
		return "Struct"
	}
	return "<unknown>"
}

type ArrayFiller interface {
	// when elem is nil, there is no more to fill
	Next() (index int, elem interface{})
	Fill()
}

type MapFiller interface {
	Next() (key interface{}, elem interface{})
	Fill()
}

type Accessor interface {
	fmt.GoStringer
	Kind() Kind
	// ptr
	PtrElem(obj interface{}) (elem interface{}, elemAccessor Accessor)
	SetPtrElem(obj interface{}, template interface{}) (elem interface{}, elemAccessor Accessor)
	// map
	Key() Accessor
	// array/map
	Elem() Accessor
	// struct
	NumField() int
	Field(index int) StructField
	// map
	IterateMap(obj interface{}, cb func(key interface{}, elem interface{}) bool)
	FillMap(obj interface{}, cb func(filler MapFiller))
	// array
	IterateArray(obj interface{}, cb func(index int, elem interface{}) bool)
	FillArray(obj interface{}, cb func(filler ArrayFiller))
	// primitives
	Skip(obj interface{})
	Int(obj interface{}) int
	SetInt(obj interface{}, val int)
	String(obj interface{}) string
	SetString(obj interface{}, val string)
	Uintptr(obj interface{}) uintptr
}

type StructField interface {
	Name() string
	Accessor() Accessor
	Tags() map[string]interface{}
}

type NoopAccessor struct {
	AccessorTypeName string
}

func (accessor *NoopAccessor) reportError() string {
	panic(fmt.Sprintf("%s: unsupported operation", accessor.AccessorTypeName))
}

func (accessor *NoopAccessor) PtrElem(obj interface{}) (elem interface{}, elemAccessor Accessor) {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) SetPtrElem(obj interface{}, template interface{}) (elem interface{}, elemAccessor Accessor) {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) Key() Accessor {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) Elem() Accessor {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) NumField() int {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) Field(index int) StructField {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) IterateMap(obj interface{}, cb func(key interface{}, elem interface{}) bool) {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) FillMap(obj interface{}, cb func(filler MapFiller)) {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) IterateArray(obj interface{}, cb func(index int, elem interface{}) bool) {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) FillArray(obj interface{}, cb func(filler ArrayFiller)) {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) Int(obj interface{}) int {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) SetInt(obj interface{}, val int) {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) String(obj interface{}) string {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) SetString(obj interface{}, val string) {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) Uintptr(obj interface{}) uintptr {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) Skip(obj interface{}) {
}
