package native

import (
	"github.com/v2pro/plz"
	"github.com/v2pro/plz/acc"
	"reflect"
	"unsafe"
	"fmt"
)

type structAccessor struct {
	acc.NoopAccessor
	typ reflect.Type
}

func (accessor *structAccessor) Kind() reflect.Kind {
	return reflect.Struct
}

func (accessor *structAccessor) GoString() string {
	return accessor.typ.Name()
}

func (accessor *structAccessor) NumField() int {
	return accessor.typ.NumField()
}

func (accessor *structAccessor) Field(index int) acc.StructField {
	field := accessor.typ.Field(index)
	ptrType := reflect.PtrTo(field.Type)
	fieldAcc := plz.AccessorOf(ptrType)
	templateObj := castToEmptyInterface(reflect.New(field.Type).Interface())
	return acc.StructField{
		Name: field.Name,
		Accessor: &structFieldAccessor{
			structAccessor: accessor,
			field:          field,
			templateObj:    templateObj,
			accessor:       fieldAcc,
		},
	}
}

type structFieldAccessor struct {
	acc.NoopAccessor
	structAccessor acc.Accessor
	field          reflect.StructField
	templateObj    emptyInterface
	accessor       acc.Accessor
}

func (acc *structFieldAccessor) Kind() reflect.Kind {
	return acc.accessor.Kind()
}

func (acc *structFieldAccessor) Int(obj interface{}) int {
	return acc.accessor.Int(acc.fieldOf(obj))
}

func (acc *structFieldAccessor) SetInt(obj interface{}, val int) {
	acc.accessor.SetInt(acc.fieldOf(obj), val)
}

func (acc *structFieldAccessor) String(obj interface{}) string {
	return acc.accessor.String(acc.fieldOf(obj))
}

func (acc *structFieldAccessor) SetString(obj interface{}, val string) {
	acc.accessor.SetString(acc.fieldOf(obj), val)
}

func (acc *structFieldAccessor) fieldOf(obj interface{}) interface{} {
	structPtr := uintptr(extractPtrFromEmptyInterface(obj))
	structFieldPtr := structPtr + acc.field.Offset
	objEmptyInterface := acc.templateObj
	objEmptyInterface.word = unsafe.Pointer(structFieldPtr)
	return castBackEmptyInterface(objEmptyInterface)
}

func (acc *structFieldAccessor) GoString() string {
	return fmt.Sprintf("%#v/%s %#v", acc.structAccessor.GoString(), acc.field.Name, acc.accessor.GoString())
}
