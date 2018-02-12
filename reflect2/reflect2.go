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

type ArrayType interface {
	Type
	Elem() Type
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
	NumField() int
	Field(i int) StructField
	FieldByName(name string) StructField
	FieldByIndex(index []int) StructField
	FieldByNameFunc(match func(string) bool) StructField
}

type StructField interface {
	Name() string
	PkgPath() string
	Type() Type
	Tag() reflect.StructTag
	Index() []int
	Anonymous() bool
	Set(obj interface{}, value interface{})
	UnsafeSet(obj unsafe.Pointer, value unsafe.Pointer)
	Get(obj interface{}) interface{}
	UnsafeGet(obj unsafe.Pointer) unsafe.Pointer
}

type MapType interface {
	Type
	Elem() Type
	MakeMap(cap int) interface{}
	UnsafeMakeMap(cap int) unsafe.Pointer
	Set(obj interface{}, key interface{}, elem interface{})
	UnsafeSet(obj unsafe.Pointer, key unsafe.Pointer, elem unsafe.Pointer)
	TryGet(obj interface{}, key interface{}) (interface{}, bool)
	Get(obj interface{}, key interface{}) interface{}
	UnsafeGet(obj unsafe.Pointer, key unsafe.Pointer) unsafe.Pointer
	Iterate(obj interface{}) MapIterator
	UnsafeIterate(obj unsafe.Pointer) MapIterator
}

type MapIterator interface {
	HasNext() bool
	Next() (key interface{}, elem interface{})
	UnsafeNext() (key unsafe.Pointer, elem unsafe.Pointer)
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

var ConfigUnsafe = Config{UseSafeImplementation: false}.Froze()
var ConfigSafe = Config{UseSafeImplementation: true}.Froze()

func (cfg *frozenConfig) TypeOf(obj interface{}) Type {
	valType := reflect.TypeOf(obj)
	return cfg.Type2(valType)
}

func (cfg *frozenConfig) Type2(type1 reflect.Type) Type {
	safeType := safeType{Type: type1, cfg: cfg}
	switch type1.Kind() {
	case reflect.Struct:
		if cfg.useSafeImplementation {
			return &safeStructType{safeType}
		}
		return newUnsafeStructType(cfg, type1)
	case reflect.Array:
		if cfg.useSafeImplementation {
			return &safeSliceType{safeType}
		}
		return newUnsafeArrayType(cfg, type1)
	case reflect.Slice:
		if cfg.useSafeImplementation {
			return &safeSliceType{safeType}
		}
		return newUnsafeSliceType(cfg, type1)
	case reflect.Map:
		if cfg.useSafeImplementation {
			return &safeMapType{safeType}
		}
		return newUnsafeMapType(cfg, type1)
	default:
		if cfg.useSafeImplementation {
			return &safeType
		}
		return newUnsafeType(cfg, type1)
	}
}

func TypeOf(obj interface{}) Type {
	return ConfigUnsafe.TypeOf(obj)
}

func Type2(type1 reflect.Type) Type {
	return ConfigUnsafe.Type2(type1)
}

func PtrOf(obj interface{}) unsafe.Pointer {
	return unpackEFace(obj).data
}
