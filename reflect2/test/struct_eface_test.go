package test

import (
	"testing"
	"github.com/v2pro/plz/reflect2"
)

func Test_struct_eface(t *testing.T) {
	type TestObject struct {
		Field1 interface{}
	}
	t.Run("Set", testOp(func(api reflect2.API) interface{} {
		valType := api.TypeOf(TestObject{}).(reflect2.StructType)
		field1 := valType.FieldByName("Field1")
		obj := TestObject{}
		field1.Set(&obj, 100)
		return obj
	}))
	t.Run("Get", testOp(func(api reflect2.API) interface{} {
		obj := TestObject{Field1: 100}
		valType := api.TypeOf(obj).(reflect2.StructType)
		field1 := valType.FieldByName("Field1")
		return field1.Get(&obj)
	}))
}