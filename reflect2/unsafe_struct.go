package reflect2

import "reflect"

type unsafeStructType struct {
	unsafeType
}

func newUnsafeStructType(cfg *frozenConfig, type1 reflect.Type) *unsafeStructType {
	return &unsafeStructType{
		unsafeType: *newUnsafeType(cfg, type1),
	}
}

func (type2 *unsafeStructType) FieldByName(name string) StructField {
	structField, found := type2.Type.FieldByName(name)
	if !found {
		panic("field " + name + " not found")
	}
	switch structField.Type.Kind() {
	case reflect.Interface:
		return &unsafeEFaceField{StructField: structField}
	default:
		return &unsafeDirField{
			StructField: structField,
			rtype:       unpackEFace(structField.Type).data,
		}
	}
}