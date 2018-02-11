package test

import (
	"testing"
	"github.com/v2pro/plz/test"
	"github.com/v2pro/plz/countlog"
	"github.com/v2pro/plz/reflect2"
	"reflect"
	"github.com/v2pro/plz/test/must"
)

func Test_level1(t *testing.T) {
	t.Run("new int", test.Case(func(ctx *countlog.Context) {
		valType := reflect2.TypeOf(1)
		obj := valType.New()
		reflect.ValueOf(obj).Elem().Set(reflect.ValueOf(100))
		must.Equal(100, *obj.(*int))
	}))
}
