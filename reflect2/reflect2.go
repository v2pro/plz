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
	// Type1 returns reflect.Type
	Type1() reflect.Type
	FieldByName(fieldName string) StructField
}

type StructField interface {
	Set(obj interface{}, value interface{})
	UnsafeSet(obj unsafe.Pointer, value unsafe.Pointer)
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
}

var ConfigUnsafe = Config{UseSafeImplementation:false}.Froze()
var ConfigSafe = Config{UseSafeImplementation:true}.Froze()


func (cfg *frozenConfig) TypeOf(obj interface{}) Type {
	valType := reflect.TypeOf(obj)
	if cfg.useSafeImplementation {
		return &safeType{Type: valType}
	}
	rtype := toEface(valType).data
	prtype := toEface(reflect.PtrTo(valType)).data
	return &unsafeType{
		Type:   valType,
		rtype:  rtype,
		prtype: prtype,
	}
}

func TypeOf(obj interface{}) Type {
	return ConfigUnsafe.TypeOf(obj)
}
