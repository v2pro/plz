package plz

import (
	"reflect"
	"github.com/v2pro/plz/acc"
)

func AccessorOf(typ reflect.Type) lang.Accessor {
	return lang.AccessorOf(typ)
}