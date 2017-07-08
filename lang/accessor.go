package lang

import (
	"fmt"
	"reflect"
	"unsafe"
	"github.com/v2pro/plz/lang/tagging"
)

var AccessorProviders = []func(typ reflect.Type, tagName string) Accessor{}

func AccessorOf(typ reflect.Type, tagName string) Accessor {
	for _, provider := range AccessorProviders {
		accessor := provider(typ, tagName)
		if accessor != nil {
			return accessor
		}
	}
	panic(fmt.Sprintf("no accessor provider for: %v", typ))
}

func AddressOf(obj interface{}) unsafe.Pointer {
	ptr := castToEmptyInterface(obj).word
	typ := reflect.TypeOf(obj)
	switch typ.Kind() {
	case reflect.Array:
		if typ.Len() == 1 && (typ.Elem().Kind() == reflect.Ptr || typ.Elem().Kind() == reflect.Map){
			asVal := uintptr(ptr)
			ptr = unsafe.Pointer(&asVal)
		}
	case reflect.Struct:
		if typ.NumField() == 1 && (typ.Field(0).Type.Kind() == reflect.Ptr || typ.Field(0).Type.Kind() == reflect.Map) {
			asVal := uintptr(ptr)
			ptr = unsafe.Pointer(&asVal)
		}
	}
	return ptr
}

func castToEmptyInterface(obj interface{}) emptyInterface {
	return *((*emptyInterface)(unsafe.Pointer(&obj)))
}

// emptyInterface is the header for an interface{} value.
type emptyInterface struct {
	typ  unsafe.Pointer
	word unsafe.Pointer
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
		return "Variant"
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
	Next() (index int, elem unsafe.Pointer)
	Fill()
}

type MapFiller interface {
	Next() (key unsafe.Pointer, elem unsafe.Pointer)
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
	New() (interface{}, Accessor)

	// === runtime ===
	IsNil(ptr unsafe.Pointer) bool
	// variant
	VariantElem(ptr unsafe.Pointer) (elem unsafe.Pointer, elemAccessor Accessor)
	InitVariant(ptr unsafe.Pointer, template Accessor) (elem unsafe.Pointer, elemAccessor Accessor)
	// map
	MapIndex(ptr unsafe.Pointer, key unsafe.Pointer) (elem unsafe.Pointer) // only when random accessible
	SetMapIndex(ptr unsafe.Pointer, key unsafe.Pointer, elem unsafe.Pointer) // only when random accessible
	IterateMap(ptr unsafe.Pointer, cb func(key unsafe.Pointer, elem unsafe.Pointer) bool)
	FillMap(ptr unsafe.Pointer, cb func(filler MapFiller))
	// array/struct
	ArrayIndex(ptr unsafe.Pointer, index int) (elem unsafe.Pointer) // only when random accessible
	IterateArray(ptr unsafe.Pointer, cb func(index int, elem unsafe.Pointer) bool)
	FillArray(ptr unsafe.Pointer, cb func(filler ArrayFiller))
	// primitives
	Skip(ptr unsafe.Pointer) // when the value is not needed
	String(ptr unsafe.Pointer) string
	SetString(ptr unsafe.Pointer, val string)
	Bool(ptr unsafe.Pointer) bool
	SetBool(ptr unsafe.Pointer, val bool)
	Int(ptr unsafe.Pointer) int
	SetInt(ptr unsafe.Pointer, val int)
	Int8(ptr unsafe.Pointer) int8
	SetInt8(ptr unsafe.Pointer, val int8)
	Int16(ptr unsafe.Pointer) int16
	SetInt16(ptr unsafe.Pointer, val int16)
	Int32(ptr unsafe.Pointer) int32
	SetInt32(ptr unsafe.Pointer, val int32)
	Int64(ptr unsafe.Pointer) int64
	SetInt64(ptr unsafe.Pointer, val int64)
	Uint(ptr unsafe.Pointer) uint
	SetUint(ptr unsafe.Pointer, val uint)
	Uint8(ptr unsafe.Pointer) uint8
	SetUint8(ptr unsafe.Pointer, val uint8)
	Uint16(ptr unsafe.Pointer) uint16
	SetUint16(ptr unsafe.Pointer, val uint16)
	Uint32(ptr unsafe.Pointer) uint32
	SetUint32(ptr unsafe.Pointer, val uint32)
	Uint64(ptr unsafe.Pointer) uint64
	SetUint64(ptr unsafe.Pointer, val uint64)
	Float32(ptr unsafe.Pointer) float32
	SetFloat32(ptr unsafe.Pointer, val float32)
	Float64(ptr unsafe.Pointer) float64
	SetFloat64(ptr unsafe.Pointer, val float64)
}

type StructField interface {
	Index() int
	Name() string
	Accessor() Accessor
	Tags() map[string]tagging.TagValue
}

type NoopAccessor struct {
	TagName string
	AccessorTypeName string
}

func (accessor *NoopAccessor) reportError() string {
	panic(fmt.Sprintf("%s: unsupported operation", accessor.AccessorTypeName))
}

func (accessor *NoopAccessor) VariantElem(ptr unsafe.Pointer) (elem unsafe.Pointer, elemAccessor Accessor) {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) InitVariant(ptr unsafe.Pointer, template Accessor) (elem unsafe.Pointer, elemAccessor Accessor) {
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

func (accessor *NoopAccessor) New() (interface{}, Accessor) {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) IsNil(ptr unsafe.Pointer) bool {
	return ptr == nil
}

func (accessor *NoopAccessor) ArrayIndex(ptr unsafe.Pointer, index int) (elem unsafe.Pointer) {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) SetMapIndex(ptr unsafe.Pointer, key unsafe.Pointer, elem unsafe.Pointer) {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) MapIndex(ptr unsafe.Pointer, key unsafe.Pointer) (elem unsafe.Pointer) {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) IterateMap(ptr unsafe.Pointer, cb func(key unsafe.Pointer, elem unsafe.Pointer) bool) {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) FillMap(ptr unsafe.Pointer, cb func(filler MapFiller)) {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) IterateArray(ptr unsafe.Pointer, cb func(index int, elem unsafe.Pointer) bool) {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) FillArray(ptr unsafe.Pointer, cb func(filler ArrayFiller)) {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) Skip(ptr unsafe.Pointer) {
}

func (accessor *NoopAccessor) String(ptr unsafe.Pointer) string {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) SetString(ptr unsafe.Pointer, val string) {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) Bool(ptr unsafe.Pointer) bool {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) SetBool(ptr unsafe.Pointer, val bool) {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) Int(ptr unsafe.Pointer) int {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) SetInt(ptr unsafe.Pointer, val int) {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) Int8(ptr unsafe.Pointer) int8 {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) SetInt8(ptr unsafe.Pointer, val int8) {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) Int16(ptr unsafe.Pointer) int16 {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) SetInt16(ptr unsafe.Pointer, val int16) {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) Int32(ptr unsafe.Pointer) int32 {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) SetInt32(ptr unsafe.Pointer, val int32) {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) Int64(ptr unsafe.Pointer) int64 {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) SetInt64(ptr unsafe.Pointer, val int64) {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) Uint(ptr unsafe.Pointer) uint {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) SetUint(ptr unsafe.Pointer, val uint) {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) Uint8(ptr unsafe.Pointer) uint8 {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) SetUint8(ptr unsafe.Pointer, val uint8) {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) Uint16(ptr unsafe.Pointer) uint16 {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) SetUint16(ptr unsafe.Pointer, val uint16) {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) Uint32(ptr unsafe.Pointer) uint32 {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) SetUint32(ptr unsafe.Pointer, val uint32) {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) Uint64(ptr unsafe.Pointer) uint64 {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) SetUint64(ptr unsafe.Pointer, val uint64) {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) Float32(ptr unsafe.Pointer) float32 {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) SetFloat32(ptr unsafe.Pointer, val float32) {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) Float64(ptr unsafe.Pointer) float64 {
	panic(accessor.reportError())
}

func (accessor *NoopAccessor) SetFloat64(ptr unsafe.Pointer, val float64) {
	panic(accessor.reportError())
}