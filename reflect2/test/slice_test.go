package test

import (
	"testing"
	"github.com/v2pro/plz/reflect2"
)

func Test_slice(t *testing.T) {
	t.Run("New", testOp(func(api reflect2.API) interface{} {
		valType := reflect2.TypeOf([]int{})
		obj := *valType.New().(*[]int)
		obj = append(obj, 1)
		return obj
	}))
	t.Run("MakeSlice", testOp(func(api reflect2.API) interface{} {
		valType := api.TypeOf([]int{}).(reflect2.SliceType)
		obj := valType.MakeSlice(5, 10)
		obj.([]int)[0] = 100
		obj.([]int)[4] = 20
		return obj
	}))
}
