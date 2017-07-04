package lang

import (
	"fmt"
	"reflect"
)

var AccessorProviders = []func(typ reflect.Type) Accessor{}

func AccessorOf(typ reflect.Type) Accessor {
	for _, provider := range AccessorProviders {
		accessor := provider(typ)
		if accessor != nil {
			return accessor
		}
	}
	panic(fmt.Sprintf("no accessor provider for: %v", typ))
}

type Kind uint

const (
	Invalid Kind = iota
	Array
	Map
	Struct
	Variant
	String
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
	Float32
	Float64
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
	case Float32:
		return true
	case Float64:
		return true
	case Array:
		return false
	case Variant:
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
	case Float32:
		return "Float32"
	case Float64:
		return "Float64"
	case Array:
		return "Array"
	case Variant:
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
	// === static ===
	fmt.GoStringer
	Kind() Kind
	// map
	Key() Accessor
	// array/map
	Elem() Accessor
	// struct
	NumField() int
	Field(index int) StructField
	// array/struct
	RandomAccessible() bool

	// === runtime ===
	// variant
	VariantElem(obj interface{}) (elem interface{}, elemAccessor Accessor)
	InitVariant(obj interface{}, template interface{}) (elem interface{}, elemAccessor Accessor)
	// map
	IterateMap(obj interface{}, cb func(key interface{}, elem interface{}) bool)
	FillMap(obj interface{}, cb func(filler MapFiller))
	// array/struct
	Index(obj interface{}, index int) (elem interface{}) // only when random accessible
	IterateArray(obj interface{}, cb func(index int, elem interface{}) bool)
	FillArray(obj interface{}, cb func(filler ArrayFiller))
	// primitives
	Skip(obj interface{}) // when the value is not needed
	String(obj interface{}) string
	SetString(obj interface{}, val string)
	Bool(obj interface{}) bool
	SetBool(obj interface{}, val bool)
	Int(obj interface{}) int
	SetInt(obj interface{}, val int)
	Int8(obj interface{}) int8
	SetInt8(obj interface{}, val int8)
	Int16(obj interface{}) int16
	SetInt16(obj interface{}, val int16)
	Int32(obj interface{}) int32
	SetInt32(obj interface{}, val int32)
	Int64(obj interface{}) int64
	SetInt64(obj interface{}, val int64)
	Uint(obj interface{}) uint
	SetUint(obj interface{}, val uint)
	Uint8(obj interface{}) uint8
	SetUint8(obj interface{}, val uint8)
	Uint16(obj interface{}) uint16
	SetUint16(obj interface{}, val uint16)
	Uint32(obj interface{}) uint32
	SetUint32(obj interface{}, val uint32)
	Uint64(obj interface{}) uint64
	SetUint64(obj interface{}, val uint64)
	Float32(obj interface{}) float32
	SetFloat32(obj interface{}, val float32)
	Float64(obj interface{}) float64
	SetFloat64(obj interface{}, val float64)
	// pointer to memory address
	AddressOf(obj interface{}) uintptr
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

func (accessor *NoopAccessor) VariantElem(obj interface{}) (elem interface{}, elemAccessor Accessor) {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) InitVariant(obj interface{}, template interface{}) (elem interface{}, elemAccessor Accessor) {
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

func (accessor *NoopAccessor) RandomAccessible() bool {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) Index(obj interface{}, index int) (elem interface{}) {
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

func (accessor *NoopAccessor) Skip(obj interface{}) {
}

func (accessor *NoopAccessor) String(obj interface{}) string {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) SetString(obj interface{}, val string) {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) Bool(obj interface{}) bool {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) SetBool(obj interface{}, val bool) {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) Int(obj interface{}) int {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) SetInt(obj interface{}, val int) {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) Int8(obj interface{}) int8 {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) SetInt8(obj interface{}, val int8) {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) Int16(obj interface{}) int16 {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) SetInt16(obj interface{}, val int16) {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) Int32(obj interface{}) int32 {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) SetInt32(obj interface{}, val int32) {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) Int64(obj interface{}) int64 {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) SetInt64(obj interface{}, val int64) {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) Uint(obj interface{}) uint {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) SetUint(obj interface{}, val uint) {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) Uint8(obj interface{}) uint8 {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) SetUint8(obj interface{}, val uint8) {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) Uint16(obj interface{}) uint16 {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) SetUint16(obj interface{}, val uint16) {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) Uint32(obj interface{}) uint32 {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) SetUint32(obj interface{}, val uint32) {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) Uint64(obj interface{}) uint64 {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) SetUint64(obj interface{}, val uint64) {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) Float32(obj interface{}) float32 {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) SetFloat32(obj interface{}, val float32) {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) Float64(obj interface{}) float64 {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) SetFloat64(obj interface{}, val float64) {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) AddressOf(obj interface{}) uintptr {
	panic(accessor.reportError())
}
