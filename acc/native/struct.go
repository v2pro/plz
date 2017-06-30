package native

import (
	"github.com/v2pro/plz"
	"github.com/v2pro/plz/acc"
	"reflect"
	"unsafe"
	"fmt"
	"github.com/v2pro/plz/tagging"
)

func accessorOfStruct(typ reflect.Type) acc.Accessor {
	tags := tagging.Get(typ)
	fields := []acc.StructField{}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		ptrType := reflect.PtrTo(field.Type)
		fieldAcc := plz.AccessorOf(ptrType)
		templateObj := castToEmptyInterface(reflect.New(field.Type).Interface())
		fieldTags := tags.Fields[field.Name]
		if fieldTags == nil {
			fieldTags = map[string]interface{}{}
		}
		fields = append(fields, acc.StructField{
			Name: field.Name,
			Tags: fieldTags,
			Accessor: &structFieldAccessor{
				structName:  typ.Name(),
				field:       field,
				templateObj: templateObj,
				accessor:    fieldAcc,
			},
		})
	}
	return &structAccessor{
		typ:    typ,
		fields: fields,
	}
}

type structAccessor struct {
	acc.NoopAccessor
	typ    reflect.Type
	fields []acc.StructField
}

func (accessor *structAccessor) Kind() acc.Kind {
	return acc.Struct
}

func (accessor *structAccessor) GoString() string {
	return accessor.typ.Name()
}

func (accessor *structAccessor) Key() acc.Accessor {
	return &stringAccessor{}
}

func (accessor *structAccessor) Elem() acc.Accessor {
	return &structValueAccessor{}
}

func (accessor *structAccessor) IterateMap(obj interface{}, cb func(key interface{}, elem interface{}) bool) {
	for i := 0; i < len(accessor.fields); i++ {
		field := accessor.fields[i]
		cb(field.Name, field.Accessor.Interface(obj))
	}
}

func (accessor *structAccessor) NumField() int {
	return len(accessor.fields)
}

func (accessor *structAccessor) Field(index int) acc.StructField {
	return accessor.fields[index]
}

type structFieldAccessor struct {
	acc.NoopAccessor
	structName  string
	field       reflect.StructField
	templateObj emptyInterface
	accessor    acc.Accessor
}

func (accessor *structFieldAccessor) Kind() acc.Kind {
	return accessor.accessor.Kind()
}

func (accessor *structFieldAccessor) NumField() int {
	return accessor.accessor.NumField()
}

func (accessor *structFieldAccessor) Field(index int) acc.StructField {
	return accessor.accessor.Field(index)
}

func (accessor *structFieldAccessor) Uintptr(obj interface{}) uintptr {
	structPtr := uintptr(extractPtrFromEmptyInterface(obj))
	structFieldPtr := structPtr + accessor.field.Offset
	return structFieldPtr
}

func (accessor *structFieldAccessor) Int(obj interface{}) int {
	return accessor.accessor.Int(accessor.fieldOf(obj))
}

func (accessor *structFieldAccessor) SetInt(obj interface{}, val int) {
	accessor.accessor.SetInt(accessor.fieldOf(obj), val)
}

func (accessor *structFieldAccessor) String(obj interface{}) string {
	return accessor.accessor.String(accessor.fieldOf(obj))
}

func (accessor *structFieldAccessor) SetString(obj interface{}, val string) {
	accessor.accessor.SetString(accessor.fieldOf(obj), val)
}

func (accessor *structFieldAccessor) Interface(obj interface{}) interface{} {
	return accessor.fieldOf(obj)
}

func (accessor *structFieldAccessor) fieldOf(obj interface{}) interface{} {
	structPtr := uintptr(extractPtrFromEmptyInterface(obj))
	structFieldPtr := structPtr + accessor.field.Offset
	objEmptyInterface := accessor.templateObj
	objEmptyInterface.word = unsafe.Pointer(structFieldPtr)
	return castBackEmptyInterface(objEmptyInterface)
}

func (accessor *structFieldAccessor) GoString() string {
	return fmt.Sprintf("%#v/%s %#v", accessor.structName, accessor.field.Name, accessor.accessor.GoString())
}

type structValueAccessor struct {
	acc.NoopAccessor
}

func (accessor *structValueAccessor) Kind() acc.Kind {
	return acc.Interface
}

func (accessor *structValueAccessor) GoString() string {
	return "interface{}"
}

func (accessor *structValueAccessor) String(obj interface{}) string {
	return *((*string)(extractPtrFromEmptyInterface(obj)))
}

func (accessor *structValueAccessor) SetString(obj interface{}, val string) {
	if reflect.TypeOf(obj).Kind() != reflect.Ptr {
		panic("can only SetString on pointer")
	}
	*((*string)(extractPtrFromEmptyInterface(obj))) = val
}

