package test

import (
	"testing"
	"github.com/v2pro/plz/reflect2"
	"reflect"
)

func Test_struct(t *testing.T) {
	t.Run("New", testOp(func(api reflect2.API) interface{} {
		type TestObject struct {
			Field1 int
			Field2 int
		}
		valType := api.TypeOf(TestObject{})
		obj := valType.New()
		reflect.ValueOf(obj).Elem().FieldByName("Field1").Set(reflect.ValueOf(20))
		reflect.ValueOf(obj).Elem().FieldByName("Field2").Set(reflect.ValueOf(100))
		return obj
	}))
	t.Run("Set", testOp(func(api reflect2.API) interface{} {
		type TestObject struct {
			Field1 int
			Field2 int
		}
		valType := api.TypeOf(TestObject{})
		field1 := valType.FieldByName("Field1")
		obj := TestObject{}
		field1.Set(&obj, 100)
		return obj
	}))
}
