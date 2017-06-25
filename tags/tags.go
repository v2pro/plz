package tags

import (
	"unsafe"
	"reflect"
	"strconv"
)

var registry = map[reflect.Type]interface{}{}

type Protocol interface {
	DefineTags() Tags
}

type Tags []interface{}

type StructTags struct {
	Struct map[string]interface{}
	Fields map[string]FieldTags
}

type FieldTags map[string]interface{}

var protocolType = reflect.TypeOf((*Protocol)(nil)).Elem()

func Get(typ reflect.Type) *StructTags {
	structTags, found := registry[typ].(*StructTags)
	if found {
		return structTags
	}
	if !reflect.PtrTo(typ).ConvertibleTo(protocolType) {
		return nil
	}
	fakeStructPtrVal := reflect.New(typ)
	proto := fakeStructPtrVal.Interface().(Protocol)
	allDef := proto.DefineTags()
	return register(typ, fakeStructPtrVal, allDef)
}

func D(kv ...interface{}) Tags {
	return kv
}

func Define(callback interface{}) {
	callbackType := reflect.TypeOf(callback)
	structPtrType := callbackType.In(0)
	if structPtrType.Kind() != reflect.Ptr {
		panic("defineTags callback parameter should be pointer")
	}
	structType := structPtrType.Elem()
	fakeStructPtrVal := reflect.New(structType)
	ret := reflect.ValueOf(callback).Call([]reflect.Value{fakeStructPtrVal})
	allDef := ret[0].Interface().(Tags)
	register(structType, fakeStructPtrVal, allDef)
}

func register(structType reflect.Type, fakeStructPtrVal reflect.Value, allDef Tags) *StructTags {
	fakeStructPtr := extractPtr(fakeStructPtrVal.Interface())
	structTags := &StructTags{
		Fields: map[string]FieldTags{},
	}
	fakeFieldsMap := map[uintptr]string{}
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fakeFieldPtr := fakeStructPtr + field.Offset
		fieldName := field.Name
		fakeFieldsMap[fakeFieldPtr] = fieldName
		structTags.Fields[fieldName] = parseFieldTag(field.Tag)
	}
	structTags.Struct = toMap(allDef[0].(Tags))
	for _, fieldDefObj := range allDef[1:] {
		fieldDef := fieldDefObj.(Tags)
		fieldPtr := extractPtr(fieldDef[0])
		fieldName := fakeFieldsMap[fieldPtr]
		if fieldName == "" {
			panic("field not found")
		}
		fieldTags := structTags.Fields[fieldName]
		for i := 1; i < len(fieldDef); i += 2 {
			fieldTags[fieldDef[i].(string)] = fieldDef[i+1]
		}
	}
	registry[structType] = structTags
	return structTags
}

func parseFieldTag(tag reflect.StructTag) FieldTags {
	fieldTags := FieldTags{}
	for tag != "" {
		// Skip leading space.
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}
		// Scan to colon. A space, a quote or a control character is a syntax error.
		// Strictly speaking, control chars include the range [0x7f, 0x9f], not just
		// [0x00, 0x1f], but in practice, we ignore the multi-byte control characters
		// as it is simpler to inspect the tag's bytes than the tag's runes.
		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := string(tag[:i+1])
		tag = tag[i+1:]

		value, err := strconv.Unquote(qvalue)
		if err != nil {
			panic(err.Error())
		}
		fieldTags[name] = value
	}
	return fieldTags
}

func toMap(tags Tags) map[string]interface{} {
	m := map[string]interface{}{}
	for i := 0; i < len(tags); i += 2 {
		m[tags[i].(string)] = tags[i+1]
	}
	return m
}

func extractPtr(val interface{}) uintptr {
	return (*((*emptyInterface)(unsafe.Pointer(&val)))).word
}

// emptyInterface is the header for an interface{} value.
type emptyInterface struct {
	typ  uintptr
	word uintptr
}
