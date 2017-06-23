package functional

import (
	"unsafe"
	"reflect"
)

type iterateElements func(col unsafe.Pointer, f func(elem unsafe.Pointer) bool)

func genIterateElements(typ reflect.Type) iterateElements {
	return func(col unsafe.Pointer, f func(elem unsafe.Pointer) bool) {
		arr := *(*[]int)(col)
		for _, elem := range arr {
			if !f(unsafe.Pointer(&elem)) {
				return
			}
		}
	}
}
