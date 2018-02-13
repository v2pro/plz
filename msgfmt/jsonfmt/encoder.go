package jsonfmt

import (
	"unsafe"
	"reflect"
	"strings"
	"unicode"
	"sync"
	"fmt"
	"encoding/json"
	"context"
	"github.com/v2pro/plz/reflect2"
)

var bytesType = reflect2.TypeOf([]byte(nil))
var errorType = reflect2.TypeOfPointer((*error)(nil)).Elem()
var jsonMarshalerType = reflect2.TypeOfPointer((*json.Marshaler)(nil)).Elem()

type Encoder interface {
	Encode(ctx context.Context, space []byte, ptr unsafe.Pointer) []byte
}

type Extension interface {
	EncoderOf(prefix string, valType reflect.Type) Encoder
}

type Config struct {
	IncludesUnexported bool
	Extensions         []Extension
}

func (cfg Config) Froze() API {
	return &frozenConfig{
		includesUnexported: cfg.IncludesUnexported,
		extensions:         cfg.Extensions,
		encoderCache:       &sync.Map{},
		mapKeyEncoderCache: &sync.Map{},
	}
}

type API interface {
	EncoderOf(valType reflect.Type) Encoder
}

type frozenConfig struct {
	includesUnexported bool
	extensions         []Extension
	encoderCache       *sync.Map
	mapKeyEncoderCache *sync.Map
}

func (cfg *frozenConfig) EncoderOf(type1 reflect.Type) Encoder {
	encoderObj, found := cfg.encoderCache.Load(type1)
	if found {
		return encoderObj.(Encoder)
	}
	valType := reflect2.Type2(type1)
	encoder := encoderOf(cfg, "", valType)
	if isOnePtr(valType) {
		encoder = &onePtrInterfaceEncoder{encoder}
	}
	cfg.encoderCache.Store(type1, encoder)
	return encoder
}

var ConfigDefault = Config{}.Froze()

func EncoderOf(valType reflect.Type) Encoder {
	return ConfigDefault.EncoderOf(valType)
}

func MarshalToString(obj interface{}) string {
	encoder := EncoderOf(reflect.TypeOf(obj))
	return string(encoder.Encode(nil, nil, PtrOf(obj)))
}

func encoderOf(cfg *frozenConfig, prefix string, valType reflect2.Type) Encoder {
	for _, extension := range cfg.extensions {
		encoder := extension.EncoderOf(prefix, valType.Type1())
		if encoder != nil {
			return encoder
		}
	}
	if bytesType == valType {
		return &bytesEncoder{}
	}
	//if valType.Implements(errorType) && valType.Kind() == reflect.Ptr {
	//	sampleObj := reflect.New(valType).Elem().Interface()
	//	return &pointerEncoder{elemEncoder: &errorEncoder{
	//		sampleInterface: *(*emptyInterface)(unsafe.Pointer(&sampleObj)),
	//	}}
	//}
	//if valType.Implements(jsonMarshalerType) && valType.Kind() != reflect.Ptr {
	//	sampleObj := reflect.New(valType).Elem().Interface()
	//	return &jsonMarshalerEncoder{
	//		sampleInterface: *(*emptyInterface)(unsafe.Pointer(&sampleObj)),
	//	}
	//}
	switch valType.Kind() {
	case reflect.Bool:
		return &boolEncoder{}
	case reflect.Int8:
		return &int8Encoder{}
	case reflect.Uint8:
		return &uint8Encoder{}
	case reflect.Int16:
		return &int16Encoder{}
	case reflect.Uint16:
		return &uint16Encoder{}
	case reflect.Int32:
		return &int32Encoder{}
	case reflect.Uint32:
		return &uint32Encoder{}
	case reflect.Int64, reflect.Int:
		return &int64Encoder{}
	case reflect.Uint64, reflect.Uint, reflect.Uintptr:
		return &uint64Encoder{}
	case reflect.Float64:
		return &lossyFloat64Encoder{}
	case reflect.Float32:
		return &lossyFloat32Encoder{}
	case reflect.String:
		return &stringEncoder{}
	case reflect.Ptr:
		pointerType := valType.(reflect2.PointerType)
		elemEncoder := encoderOf(cfg, prefix+" [ptrElem]", pointerType.Elem())
		return &pointerEncoder{elemEncoder: elemEncoder}
	case reflect.Slice:
		sliceType := valType.(reflect2.SliceType)
		elemEncoder := encoderOf(cfg, prefix+" [sliceElem]", sliceType.Elem())
		return &sliceEncoder{
			elemEncoder: elemEncoder,
			sliceType:   sliceType.(*reflect2.UnsafeSliceType),
		}
	case reflect.Array:
		arrayType := valType.(reflect2.ArrayType)
		elemEncoder := encoderOf(cfg, prefix+" [sliceElem]", arrayType.Elem())
		return &arrayEncoder{
			elemEncoder: elemEncoder,
			arrayType:   arrayType.(*reflect2.UnsafeArrayType),
		}
	case reflect.Struct:
		structType := valType.(reflect2.StructType)
		return encoderOfStruct(cfg, prefix, structType)
	case reflect.Map:
		mapType := valType.(reflect2.MapType)
		return encoderOfMap(cfg, prefix, mapType)
	case reflect.Interface:
		interfaceType := valType.(reflect2.InterfaceType)
		if interfaceType.NumMethod() != 0 {
			return &ifaceEncoder{}
		}
		return &efaceEncoder{}
	}
	return &unsupportedEncoder{fmt.Sprintf(`"can not encode %s %s to json"`, valType.String(), prefix)}
}

type unsupportedEncoder struct {
	msg string
}

func (encoder *unsupportedEncoder) Encode(ctx context.Context, space []byte, ptr unsafe.Pointer) []byte {
	return append(space, encoder.msg...)
}

func encoderOfMap(cfg *frozenConfig, prefix string, valType reflect2.MapType) *mapEncoder {
	keyEncoder := encoderOfMapKey(cfg, prefix, valType.Key())
	elemType := valType.Elem()
	elemEncoder := encoderOf(cfg, prefix+" [mapElem]", elemType)
	return &mapEncoder{
		keyEncoder:  keyEncoder,
		elemEncoder: elemEncoder,
		mapType:     valType.(*reflect2.UnsafeMapType),
	}
}

func encoderOfMapKey(cfg *frozenConfig, prefix string, keyType reflect2.Type) Encoder {
	encoderObj, found := cfg.mapKeyEncoderCache.Load(keyType)
	if found {
		return encoderObj.(Encoder)
	}
	encoder := _encoderOfMapKey(cfg, prefix, keyType)
	cfg.mapKeyEncoderCache.Store(keyType, encoder)
	return encoder
}

func _encoderOfMapKey(cfg *frozenConfig, prefix string, keyType reflect2.Type) Encoder {
	keyEncoder := encoderOf(cfg, prefix+" [mapKey]", keyType)
	if keyType.Kind() == reflect.String || keyType == bytesType {
		return &mapStringKeyEncoder{keyEncoder}
	}
	if keyType.Kind() == reflect.Interface {
		return &mapInterfaceKeyEncoder{cfg: cfg, prefix: prefix}
	}
	return &mapNumberKeyEncoder{keyEncoder}
}

func isOnePtr(valType reflect2.Type) bool {
	if reflect2.IsPointerKind(valType.Kind()) {
		return true
	}
	if valType.Kind() == reflect.Struct {
		structType := valType.(reflect2.StructType)
		onlyFieldKind := structType.Field(0).Type().Kind()
		if structType.NumField() == 1 && reflect2.IsPointerKind(onlyFieldKind) {
			return true
		}
	}
	if valType.Kind() == reflect.Array {
		arrayType := valType.(reflect2.ArrayType)
		onlyElemKind := arrayType.Elem().Kind()
		if arrayType.Len() == 1 && reflect2.IsPointerKind(onlyElemKind) {
			return true
		}
	}
	return false
}

func encoderOfStruct(cfg *frozenConfig, prefix string, valType reflect2.StructType) *structEncoder {
	var fields []structEncoderField
	for i := 0; i < valType.NumField(); i++ {
		field := valType.Field(i)
		name := getFieldName(cfg, field)
		if name == "" {
			continue
		}
		prefix := ""
		if len(fields) != 0 {
			prefix += ","
		}
		prefix += `"`
		prefix += name
		prefix += `":`
		fields = append(fields, structEncoderField{
			structField: field.(*reflect2.UnsafeStructField),
			prefix:      prefix,
			encoder:     encoderOf(cfg, prefix+" ."+name, field.Type()),
		})
	}
	return &structEncoder{
		fields: fields,
	}
}

func getFieldName(cfg *frozenConfig, field reflect2.StructField) string {
	if !cfg.includesUnexported && !unicode.IsUpper(rune(field.Name()[0])) {
		return ""
	}
	if field.Type().Kind() == reflect.Func {
		return ""
	}
	if field.Type().Kind() == reflect.Chan {
		return ""
	}
	jsonTag := field.Tag().Get("json")
	if jsonTag == "" {
		return field.Name()
	}
	parts := strings.Split(jsonTag, ",")
	if parts[0] == "-" {
		return ""
	}
	if parts[0] == "" {
		return field.Name()
	}
	return parts[0]
}

func PtrOf(val interface{}) unsafe.Pointer {
	return (*emptyInterface)(unsafe.Pointer(&val)).word
}

// emptyInterface is the header for an interface{} value.
type emptyInterface struct {
	typ  unsafe.Pointer
	word unsafe.Pointer
}
