package reflect2

import "reflect"

type unsafeStructType struct {
	unsafeType
}

func newUnsafeStructType(type1 reflect.Type) *unsafeStructType {
	return &unsafeStructType{
		unsafeType: *newUnsafeType(type1),
	}
}

func (type2 *unsafeStructType) FieldByName(name string) StructField {
	structField1, found := type2.Type.FieldByName(name)
	if !found {
		panic("field " + name + " not found")
	}
	return &unsafeField{
		StructField: structField1,
		rtype:       toEface(structField1.Type).data,
	}
}