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
		return nil
	}
	return newUnsafeStructField(type2, structField)
}

func (type2 *UnsafeStructType) Field(i int) StructField {
	return newUnsafeStructField(type2, type2.Type.Field(i))
}

func (type2 *UnsafeStructType) FieldByIndex(index []int) StructField {
	return newUnsafeStructField(type2, type2.Type.FieldByIndex(index))
}

func (type2 *UnsafeStructType) FieldByNameFunc(match func(string) bool) StructField {
	structField, found := type2.Type.FieldByNameFunc(match)
	if !found {
		panic("field match condition not found in " + type2.Type.String())
	}
	return newUnsafeStructField(type2, structField)
}