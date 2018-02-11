package reflect2

import (
	"unsafe"
	"reflect"
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

type unsafeType struct {
	reflect.Type
	rtype  unsafe.Pointer
	prtype unsafe.Pointer
}

func (type2 unsafeType) Type1() reflect.Type {
	return type2.Type
}

func (type2 unsafeType) UnsafeNew() unsafe.Pointer {
	return unsafe_New(type2.rtype)
}

func (type2 unsafeType) New() interface{} {
	return packEface(type2.prtype, type2.UnsafeNew())
}

func (type2 unsafeType) FieldByName(name string) StructField {
	structField1, found := type2.Type.FieldByName(name)
	if !found {
		panic("field " + name + " not found")
	}
	return &unsafeField{
		StructField: structField1,
		rtype:       toEface(structField1.Type).data,
	}
}