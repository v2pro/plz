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

var ConfigUnsafe = Config{UseSafeImplementation: false}.Froze()
var ConfigSafe = Config{UseSafeImplementation: true}.Froze()

func (cfg *frozenConfig) TypeOf(obj interface{}) Type {
	valType := reflect.TypeOf(obj)
	return cfg.Type2(valType)
}

func (cfg *frozenConfig) Type2(type1 reflect.Type) Type {
	if cfg.useSafeImplementation {
		safeType := safeType{Type: type1, cfg: cfg}
		switch type1.Kind() {
		case reflect.Map:
			return &safeMapType{safeType}
		case reflect.Ptr:
			return &safePtrType{safeType}
		case reflect.Struct:
			return &safeStructType{safeType}
		case reflect.Slice:
			return &safeSliceType{safeType: safeType}
		}
		return &safeType
	}
	switch type1.Kind() {
	case reflect.Int:
		return newUnsafeType(cfg, type1)
	case reflect.Struct:
		return newUnsafeStructType(cfg, type1)
	case reflect.Array:
		return newUnsafeArrayType(cfg, type1)
	case reflect.Slice:
		return newUnsafeSliceType(cfg, type1)
	case reflect.Map:
		return newUnsafeMapType(cfg, type1)
	case reflect.Ptr:
		return newUnsafePointerType(cfg, type1)
	}
	panic("unsupported type: " + type1.String())
}

func TypeOf(obj interface{}) Type {
	return ConfigUnsafe.TypeOf(obj)
}

func PtrOf(obj interface{}) unsafe.Pointer {
	return unpackEFace(obj).data
}
