package plz

import (
	"reflect"
	"github.com/v2pro/plz/acc"
)

func AccessorOf(typ reflect.Type) acc.Accessor {
	return acc.AccessorOf(typ)
}
