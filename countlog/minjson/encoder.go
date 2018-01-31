package minjson

import (
	"unsafe"
	"reflect"
	"strings"
	"unicode"
)

var bytesType = reflect.TypeOf([]byte(nil))

type Encoder interface {
	Encode(space []byte, ptr unsafe.Pointer) []byte
}

func EncoderOf(valType reflect.Type) Encoder {
	return encoderOf("", valType)
}

func encoderOf(prefix string, valType reflect.Type) Encoder {
	if bytesType == valType {
		return &bytesEncoder{}
	}
	switch valType.Kind() {
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
	case reflect.Uint64, reflect.Uint:
		return &uint64Encoder{}
	case reflect.Float64:
		return &lossyFloat64Encoder{}
	case reflect.Float32:
		return &lossyFloat32Encoder{}
	case reflect.String:
		return &stringEncoder{}
	case reflect.Ptr:
		elemEncoder := encoderOf(prefix+" [ptrElem]", valType.Elem())
		return &pointerEncoder{elemEncoder: elemEncoder}
	case reflect.Slice:
		elemEncoder := encoderOf(prefix+" [sliceElem]", valType.Elem())
		return &sliceEncoder{
			elemEncoder: elemEncoder,
			elemSize:    valType.Elem().Size(),
		}
	case reflect.Array:
		elemEncoder := encoderOf(prefix+" [sliceElem]", valType.Elem())
		return &arrayEncoder{
			elemEncoder: elemEncoder,
			elemSize:    valType.Elem().Size(),
			length:      valType.Len(),
		}
	case reflect.Struct:
		return encoderOfStruct(prefix, valType)
	case reflect.Map:
		return encoderOfMap(prefix, valType)
	}
	return nil
}

func encoderOfMap(prefix string, valType reflect.Type) *mapEncoder {
	keyType := valType.Key()
	keyEncoder := encoderOf(prefix+" [mapKey]", keyType)
	hasQuote := keyType.Kind() == reflect.String || keyType == bytesType
	if !hasQuote {
		keyEncoder = &mapKeyEncoder{keyEncoder}
	}
	sampleObj := reflect.MakeMap(valType).Interface()
	elemType := valType.Elem()
	elemEncoder := encoderOf(prefix+" [mapElem]", elemType)
	if isOnePtr(elemType) {
		elemEncoder = &onePtrInterfaceEncoder{elemEncoder}
	}
	return &mapEncoder{
		keyEncoder:      keyEncoder,
		elemEncoder:     elemEncoder,
		sampleInterface: *(*emptyInterface)(unsafe.Pointer(&sampleObj)),
	}
}

func isOnePtr(valType reflect.Type) bool {
	if valType.Kind() == reflect.Ptr {
		return true
	}
	if valType.Kind() == reflect.Struct &&
		valType.NumField() == 1 &&
		valType.Field(0).Type.Kind() == reflect.Ptr {
		return true
	}
	if valType.Kind() == reflect.Array &&
		valType.Len() == 1 &&
		valType.Elem().Kind() == reflect.Ptr {
		return true
	}
	return false
}

func encoderOfStruct(prefix string, valType reflect.Type) *structEncoder {
	var fields []structEncoderField
	for i := 0; i < valType.NumField(); i++ {
		field := valType.Field(i)
		name := getFieldName(field)
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
			offset:  field.Offset,
			prefix:  prefix,
			encoder: encoderOf(prefix+" ."+name, field.Type),
		})
	}
	return &structEncoder{
		fields: fields,
	}
}

func getFieldName(field reflect.StructField) string {
	if !unicode.IsUpper(rune(field.Name[0])) {
		return ""
	}
	jsonTag := field.Tag.Get("json")
	if jsonTag == "" {
		return field.Name
	}
	parts := strings.Split(jsonTag, ",")
	if parts[0] == "-" {
		return ""
	}
	if parts[0] == "" {
		return field.Name
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
