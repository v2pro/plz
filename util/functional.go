package util

import (
	"reflect"
)

// TODO: make this thread safe
var maxSimpleValue = map[reflect.Type]func(collection []interface{}) interface{}{}
var GenMaxSimpleValue = func(typ reflect.Type) func(collection []interface{}) interface{} {
	panic("not implemented")
}

type structAndField struct {
	S reflect.Type
	F string
}

// TODO: make this thread safe
var maxStructByField = map[structAndField]func(collection []interface{}) interface{}{}
var GenMaxStructByField = func(typ reflect.Type, fieldName string) func(collection []interface{}) interface{} {
	panic("not implemented")
}

// Max takes ints/floats as input. or only one array/slice as input.
// When input is only one array/slice, the element can be interface{}/ints/floats/pointers to numbers.
// Arguments must be same type. The return type will keep original element type.
func Max(collection ...interface{}) interface{} {
	if len(collection) == 0 {
		return nil
	}
	typ := reflect.TypeOf(collection[0])
	f := maxSimpleValue[typ]
	if f != nil {
		return f(collection)
	}
	f = GenMaxSimpleValue(typ)
	if f != nil {
		maxSimpleValue[typ] = f
		return f(collection)
	}
	lastElem, isString := collection[len(collection)-1].(string)
	if isString {
		cacheKey := structAndField{typ, lastElem}
		f = maxStructByField[cacheKey]
		if f != nil {
			return f(collection[:len(collection)-1])
		}
		f = GenMaxStructByField(typ, lastElem)
		if f != nil {
			maxStructByField[cacheKey] = f
			return f(collection[:len(collection)-1])
		}
	}
	panic("not implemented")
}

// Min takes ints/floats as input. or only one array/slice as input.
// When input is only one array/slice, the element can be interface{}/ints/floats/pointers to numbers.
// Arguments must be same type. The return type will keep original element type.
var Min = func(collection ...interface{}) interface{} {
	panic("not implemented")
}
