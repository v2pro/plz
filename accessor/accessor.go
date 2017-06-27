package accessor

import "reflect"

var Providers = []func(interface{}) Accessor{}

func Of(obj interface{}) Accessor {
	for _, provider := range Providers {
		asor := provider(obj)
		if asor != nil {
			return asor
		}
	}
	panic("no accessor provider for this object")
}

type Accessor interface {
	// Kind returns the specific kind of this type.
	Kind() reflect.Kind
	Int(obj interface{}) int
	SetInt(obj interface{}, val int)
}
