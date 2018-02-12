package reflect2

import (
	"reflect"
)

type UnsafeStructType struct {
	unsafeType
}

func newUnsafeStructType(cfg *frozenConfig, type1 reflect.Type) *UnsafeStructType {
	return &UnsafeStructType{
		unsafeType: *newUnsafeType(cfg, type1),
	}
}

func (type2 *UnsafeStructType) FieldByName(name string) StructField {
	structField, found := type2.Type.FieldByName(name)
	if !found {
		panic("field " + name + " not found in " + type2.Type.String())
	}
	return &UnsafeStructField{
		StructField: structField,
		rtype:       unpackEFace(structField.Type).data,
		ptrRType:    unpackEFace(reflect.PtrTo(structField.Type)).data,
		structType:  type2,
	}
}
