package native

import (
	"github.com/v2pro/plz/acc"
	"reflect"
	"github.com/v2pro/plz/tagging"
	"unsafe"
)

func accessorOfStruct(typ reflect.Type) acc.Accessor {
	tags := tagging.Get(typ)
	fields := []*structField{}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldAcc := acc.AccessorOf(reflect.PtrTo(field.Type))
		fieldTags := tags.Fields[field.Name]
		if fieldTags == nil {
			fieldTags = map[string]interface{}{}
		}
		templateElemObj := castToEmptyInterface(reflect.New(field.Type).Interface())
		fields = append(fields, &structField{
			name:            field.Name,
			tags:            fieldTags,
			accessor:        fieldAcc,
			size:            field.Type.Size(),
			templateElemObj: templateElemObj,
		})
	}
	return &structAccessor{
		NoopAccessor: acc.NoopAccessor{"structAccessor"},
		typ:          typ,
		fields:       fields,
	}
}

type structField struct {
	name            string
	accessor        acc.Accessor
	tags            map[string]interface{}
	size            uintptr
	templateElemObj emptyInterface
}

func (sf *structField) Name() string {
	return sf.name
}

func (sf *structField) Accessor() acc.Accessor {
	return sf.accessor
}

func (sf *structField) Tags() map[string]interface{} {
	return sf.tags
}

type structAccessor struct {
	acc.NoopAccessor
	typ    reflect.Type
	fields []*structField
}

func (accessor *structAccessor) Kind() acc.Kind {
	return acc.Struct
}

func (accessor *structAccessor) GoString() string {
	return accessor.typ.Name()
}

func (accessor *structAccessor) NumField() int {
	return len(accessor.fields)
}

func (accessor *structAccessor) Field(index int) acc.StructField {
	return accessor.fields[index]
}

func (accessor *structAccessor) IterateArray(obj interface{}, cb func(index int, elem interface{}) bool) {
	currentPtr := uintptr(extractPtrFromEmptyInterface(obj))
	for index := 0; index < len(accessor.fields); index++ {
		field := accessor.fields[index]
		currentObj := field.templateElemObj
		currentObj.word = unsafe.Pointer(currentPtr)
		cb(index, castBackEmptyInterface(currentObj))
		currentPtr += field.size
	}
}

func (accessor *structAccessor) FillArray(obj interface{}, cb func(filler acc.ArrayFiller)) {
	filler := &structFiller{
		fields:     accessor.fields,
		currentPtr: uintptr(extractPtrFromEmptyInterface(obj)),
	}
	cb(filler)
}

type structFiller struct {
	fields     []*structField
	index      int
	currentPtr uintptr
}

func (filler *structFiller) Next() (int, interface{}) {
	field := filler.fields[filler.index]
	currentObj := field.templateElemObj
	currentObj.word = unsafe.Pointer(filler.currentPtr)
	filler.currentPtr += field.size
	currentIndex := filler.index
	filler.index++
	return currentIndex, castBackEmptyInterface(currentObj)
}

func (filler *structFiller) Fill() {
}
