package reflect2

import (
	"reflect"
	"unsafe"
)

type Type interface {
	// New return pointer to data of this type
	New() interface{}
	// UnsafeNew return the allocated space pointed by unsafe.Pointer
	UnsafeNew() unsafe.Pointer
	// PackEFace cast a pointer back to empty interface
	PackEFace(ptr unsafe.Pointer) interface{}
	// Type1 returns reflect.Type
	Type1() reflect.Type
}

type StringType interface {
	Type
}

type ArrayType interface {
	Type
	Set(obj interface{}, index int, elem interface{})
	UnsafeSet(obj unsafe.Pointer, index int, elem unsafe.Pointer)
	Get(obj interface{}, index int) interface{}
	UnsafeGet(obj unsafe.Pointer, index int) unsafe.Pointer
}

type SliceType interface {
	ArrayType
	MakeSlice(length int, cap int) interface{}
	UnsafeMakeSlice(length int, cap int) unsafe.Pointer
	Append(obj interface{}, elem interface{}) interface{}
	UnsafeAppend(obj unsafe.Pointer, elem unsafe.Pointer) unsafe.Pointer
}

type StructType interface {
	Type
	FieldByName(name string) StructField
}

type StructField interface {
	Set(obj interface{}, value interface{})
	UnsafeSet(obj unsafe.Pointer, value unsafe.Pointer)
}

type MapType interface {
	Type
	MakeMap(cap int) interface{}
	UnsafeMakeMap(cap int) unsafe.Pointer
	Set(obj interface{}, key interface{}, elem interface{})
	UnsafeSet(obj unsafe.Pointer, key unsafe.Pointer, elem unsafe.Pointer)
	Iterate(obj interface{}) MapIterator
	UnsafeIterate(obj unsafe.Pointer) *UnsafeMapIterator
}

type MapIterator interface {
	HasNext() bool
	Next() (key interface{}, elem interface{})
}

type PointerType interface {
	Type
	Get(obj interface{}) interface{}
	UnsafeGet(obj unsafe.Pointer) unsafe.Pointer
}

type Config struct {
	UseSafeImplementation bool
}

func (cfg Config) Froze() *frozenConfig {
	return &frozenConfig{useSafeImplementation: cfg.UseSafeImplementation}
}

type frozenConfig struct {
	useSafeImplementation bool
}

type API interface {
	TypeOf(obj interface{}) Type
	Type2(type1 reflect.Type) Type
}

var ConfigUnsafe = Config{UseSafeImplementation:false}.Froze()
var ConfigSafe = Config{UseSafeImplementation:true}.Froze()

func (cfg *frozenConfig) TypeOf(obj interface{}) Type {
	valType := reflect.TypeOf(obj)
	return cfg.Type2(valType)
}

func (cfg *frozenConfig) Type2(type1 reflect.Type) Type {
	if cfg.useSafeImplementation {
		switch type1.Kind() {
		case reflect.Map:
			return &safeMapType{safeType{Type: type1}}
		case reflect.Ptr:
			return &safePtrType{safeType{Type: type1}}
		}
		return &safeType{Type: type1}
	}
	switch type1.Kind() {
	case reflect.Int:
		return newUnsafeType(type1)
	case reflect.Struct:
		return newUnsafeStructType(type1)
	case reflect.Array:
		return newUnsafeArrayType(type1)
	case reflect.Slice:
		return newUnsafeSliceType(type1)
	case reflect.Map:
		return newUnsafeMapType(type1)
	case reflect.Ptr:
		return newUnsafePointerType(type1)
	}
	panic("unsupported type: " + type1.String())
}

func TypeOf(obj interface{}) Type {
	return ConfigUnsafe.TypeOf(obj)
}
