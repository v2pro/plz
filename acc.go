package plz

import (
	"reflect"
	"github.com/v2pro/plz/acc"
)

func AccessorOf(dstType reflect.Type, srcType reflect.Type) acc.Accessor {
	return acc.AccessorOf(dstType, srcType)
}