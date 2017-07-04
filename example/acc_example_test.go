package example

import (
	_ "github.com/v2pro/plz/lang/nativeacc"
	"reflect"
	"fmt"
	"github.com/v2pro/plz/lang"
)

func Example_accessor() {
	obj := []int{1, 2, 3}
	accessor := lang.AccessorOf(reflect.TypeOf(obj))
	elemAccessor := accessor.Elem()
	accessor.IterateArray(obj, func(index int, elem interface{}) bool {
		fmt.Println(elemAccessor.Int(elem))
		return true
	})
	// Output:
	// 1
	// 2
	// 3
}
