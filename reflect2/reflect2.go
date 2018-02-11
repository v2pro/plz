package reflect2

import (
	"reflect"
	"unsafe"
)

//go:linkname unsafe_New reflect.unsafe_New
func unsafe_New(rtype unsafe.Pointer) unsafe.Pointer

//go:linkname typedmemmove reflect.typedmemmove
func typedmemmove(rtype unsafe.Pointer, dst, src unsafe.Pointer)

// add returns p+x.
//
// The whySafe string is ignored, so that the function still inlines
// as efficiently as p+x, but all call sites should use the string to
// record why the addition is safe, which is to say why the addition
// does not cause x to advance to the very end of p's allocation
// and therefore point incorrectly at the next block in memory.
func add(p unsafe.Pointer, x uintptr, whySafe string) unsafe.Pointer {
	return unsafe.Pointer(uintptr(p) + x)
}

type Type interface {
	// New return pointer to data of this type
	New() interface{}
	// UnsafeNew return the allocated space pointed by unsafe.Pointer
	UnsafeNew() unsafe.Pointer
	// Type1 returns reflect.Type
	Type1() reflect.Type
	FieldByName(fieldName string) *StructField
}

type StructField struct {
	reflect.StructField
	rtype unsafe.Pointer
}

func (field *StructField) Set(obj interface{}, value interface{}) {
	field.UnsafeSet(toEface(obj).data, toEface(value).data)
}

func (field *StructField) UnsafeSet(obj unsafe.Pointer, value unsafe.Pointer) {
	fieldPtr := add(obj, field.Offset, "same as non-reflect &v.field")
	typedmemmove(field.rtype, fieldPtr, value)
}

type type2 struct {
	reflect.Type
	rtype  unsafe.Pointer
	prtype unsafe.Pointer
}

func (type2 type2) Type1() reflect.Type {
	return type2.Type
}

func (type2 type2) UnsafeNew() unsafe.Pointer {
	return unsafe_New(type2.rtype)
}

func (type2 type2) New() interface{} {
	return packEface(type2.prtype, type2.UnsafeNew())
}

func (type2 type2) FieldByName(name string) *StructField {
	structField1, found := type2.Type.FieldByName(name)
	if !found {
		panic("field " + name + " not found")
	}
	return &StructField{
		StructField: structField1,
		rtype:       toEface(structField1.Type).data,
	}
}

func TypeOf(obj interface{}) Type {
	valType := reflect.TypeOf(obj)
	rtype := toEface(valType).data
	prtype := toEface(reflect.PtrTo(valType)).data
	return &type2{
		Type:   valType,
		rtype:  rtype,
		prtype: prtype,
	}
}

type eface struct {
	rtype unsafe.Pointer
	data  unsafe.Pointer
}

func toEface(obj interface{}) *eface {
	return (*eface)(unsafe.Pointer(&obj))
}

func packEface(rtype unsafe.Pointer, data unsafe.Pointer) interface{} {
	var i interface{}
	e := (*eface)(unsafe.Pointer(&i))
	e.rtype = rtype
	e.data = data
	return i
}
