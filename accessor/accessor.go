package accessor

import (
	"reflect"
	"fmt"
)

var Providers = []func(reflect.Type) Accessor{}

func Of(typ reflect.Type) Accessor {
	for _, provider := range Providers {
		asor := provider(typ)
		if asor != nil {
			return asor
		}
	}
	panic(fmt.Sprintf("no accessor provider for: %v", typ))
}

type Accessor interface {
	// Kind returns the specific kind of this type.
	Kind() reflect.Kind
	Int(obj interface{}) int
	SetInt(obj interface{}, val int)
	NumField() int
	Field(index int) StructField
}

type StructField struct {
	Name     string
	Accessor Accessor
}
