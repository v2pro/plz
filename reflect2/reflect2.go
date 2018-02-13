package reflect2

import (
	"reflect"
	"unsafe"
	"sync"
)

type Type interface {
	Kind() reflect.Kind
	// New return pointer to data of this type
	New() interface{}
	// UnsafeNew return the allocated space pointed by unsafe.Pointer
	UnsafeNew() unsafe.Pointer
	// PackEFace cast a unsafe pointer to object represented pointer
	PackEFace(ptr unsafe.Pointer) interface{}
	// Indirect dereference object represented pointer to this type
	Indirect(obj interface{}) interface{}
	// UnsafeIndirect dereference pointer to this type
	UnsafeIndirect(ptr unsafe.Pointer) interface{}
	// Type1 returns reflect.Type
	Type1() reflect.Type
	Implements(thatType Type) bool
	String() string
	RType() uintptr
}

type ListType interface {
	Type
	Elem() Type
	Set(obj interface{}, index int, elem interface{})
	UnsafeSet(obj unsafe.Pointer, index int, elem unsafe.Pointer)
	Get(obj interface{}, index int) interface{}
	UnsafeGet(obj unsafe.Pointer, index int) unsafe.Pointer
}

type ArrayType interface {
	ListType
	Len() int
}

type SliceType interface {
	ListType
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
	Key() Type
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

type PtrType interface {
	Type
	Elem() Type
}

type InterfaceType interface {
	NumMethod() int
}

type Config struct {
	UseSafeImplementation bool
}

type frozenConfig struct {
	useSafeImplementation bool
	cache *sync.Map
}

type API interface {
	TypeOf(obj interface{}) Type
	Type2(type1 reflect.Type) Type
}

var ConfigUnsafe = Config{UseSafeImplementation: false}.Froze()
var ConfigSafe = Config{UseSafeImplementation: true}.Froze()

func (cfg Config) Froze() *frozenConfig {
	return &frozenConfig{
		useSafeImplementation: cfg.UseSafeImplementation,
		cache: &sync.Map{},
	}
}

func (cfg *frozenConfig) TypeOf(obj interface{}) Type {
	cacheKey := uintptr(unpackEFace(obj).rtype)
	typeObj, found := cfg.cache.Load(cacheKey)
	if found {
		return typeObj.(Type)
	}
	return cfg.Type2(reflect.TypeOf(obj))
}

func (cfg *frozenConfig) Type2(type1 reflect.Type) Type {
	cacheKey := uintptr(unpackEFace(type1).data)
	typeObj, found := cfg.cache.Load(cacheKey)
	if found {
		return typeObj.(Type)
	}
	type2 := cfg.wrapType(type1)
	cfg.cache.Store(cacheKey, type2)
	return type2
}

func (cfg *frozenConfig) wrapType(type1 reflect.Type) Type {
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
	case reflect.Ptr, reflect.Chan, reflect.Func:
		if cfg.useSafeImplementation {
			return &safeMapType{safeType}
		}
		return newUnsafePtrType(cfg, type1)
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

func TypeOfPtr(obj interface{}) PtrType {
	return TypeOf(obj).(PtrType)
}

func Type2(type1 reflect.Type) Type {
	return ConfigUnsafe.Type2(type1)
}

func PtrOf(obj interface{}) unsafe.Pointer {
	return unpackEFace(obj).data
}

func RTypeOf(obj interface{}) uintptr {
	return uintptr(unpackEFace(obj).rtype)
}

func IsPtrKind(kind reflect.Kind) bool {
	switch kind {
	case reflect.Ptr, reflect.Map, reflect.Chan, reflect.Func:
		return true
	}
	return false
}