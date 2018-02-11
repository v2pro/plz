package reflect2

import (
	"unsafe"
	"reflect"
)

//go:linkname unsafe_New reflect.unsafe_New
func unsafe_New(rtype unsafe.Pointer) unsafe.Pointer

//go:linkname typedmemmove reflect.typedmemmove
func typedmemmove(rtype unsafe.Pointer, dst, src unsafe.Pointer)

//go:linkname unsafe_NewArray reflect.unsafe_NewArray
func unsafe_NewArray(rtype unsafe.Pointer, length int) unsafe.Pointer

// typedslicecopy copies a slice of elemType values from src to dst,
// returning the number of elements copied.
//go:linkname typedslicecopy reflect.typedslicecopy
//go:noescape
func typedslicecopy(elemType unsafe.Pointer, dst, src sliceHeader) int

//go:linkname makemap reflect.makemap
func makemap(rtype unsafe.Pointer, cap int) (m unsafe.Pointer)

//go:linkname mapassign reflect.mapassign
//go:noescape
func mapassign(rtype unsafe.Pointer, m unsafe.Pointer, key, val unsafe.Pointer)

type unsafeType struct {
	reflect.Type
	rtype  unsafe.Pointer
	ptrRType unsafe.Pointer
}

func newUnsafeType(type1 reflect.Type) *unsafeType {
	return &unsafeType{
		Type: type1,
		rtype: toEface(type1).data,
		ptrRType: toEface(reflect.PtrTo(type1)).data,
	}
}

func (type2 *unsafeType) Type1() reflect.Type {
	return type2.Type
}

func (type2 *unsafeType) UnsafeNew() unsafe.Pointer {
	return unsafe_New(type2.rtype)
}

func (type2 *unsafeType) New() interface{} {
	return packEface(type2.ptrRType, type2.UnsafeNew())
}