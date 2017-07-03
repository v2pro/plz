package plz

import (
	"reflect"
	"github.com/v2pro/plz/lang"
)

func AccessorOf(typ reflect.Type) lang.Accessor {
	return lang.AccessorOf(typ)
}

func AccessorOfObject(obj interface{}) lang.Accessor {
	return AccessorOf(reflect.TypeOf(obj))
}
