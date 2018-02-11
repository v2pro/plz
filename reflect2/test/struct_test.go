package test

import (
	"testing"
	"github.com/v2pro/plz/countlog"
	"github.com/v2pro/plz/reflect2"
	"reflect"
	"github.com/v2pro/plz/test/must"
	"github.com/v2pro/plz/test"
	"unsafe"
)

func Test_struct(t *testing.T) {
	t.Run("New", test.Case(func(ctx *countlog.Context) {
		type TestObject struct {
			Field1 int
			Field2 int
		}
		valType := reflect2.TypeOf(TestObject{})
		obj := valType.New()
		reflect.ValueOf(obj).Elem().FieldByName("Field1").Set(reflect.ValueOf(20))
		reflect.ValueOf(obj).Elem().FieldByName("Field2").Set(reflect.ValueOf(100))
		must.Equal(20, obj.(*TestObject).Field1)
		must.Equal(100, obj.(*TestObject).Field2)
	}))
	t.Run("Set", test.Case(func(ctx *countlog.Context) {
		type TestObject struct {
			Field1 int
			Field2 int
		}
		valType := reflect2.TypeOf(TestObject{})
		field1 := valType.FieldByName("Field1")
		obj := TestObject{}
		field1.Set(&obj, 100)
		must.Equal(100, obj.Field1)
	}))
	t.Run("UnsafeSet", test.Case(func(ctx *countlog.Context) {
		type TestObject struct {
			Field1 int
			Field2 int
		}
		valType := reflect2.TypeOf(TestObject{})
		field1 := valType.FieldByName("Field1")
		obj := TestObject{}
		value := 100
		field1.UnsafeSet(unsafe.Pointer(&obj), unsafe.Pointer(&value))
		must.Equal(100, obj.Field1)
	}))
}
