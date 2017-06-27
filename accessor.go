package plz

import (
	"github.com/v2pro/plz/accessor"
	"reflect"
)

func AccessorOf(typ reflect.Type) accessor.Accessor {
	return accessor.Of(typ)
}
