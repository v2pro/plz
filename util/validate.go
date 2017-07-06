package util

import (
	"github.com/v2pro/plz/lang"
	"fmt"
	"errors"
	"reflect"
	"unsafe"
)

var ValidatorProviders = []func(accessor lang.Accessor) (Validator, error){}

func Validate(obj interface{}) error {
	accessor := lang.AccessorOf(reflect.TypeOf(obj))
	validator, err := getValidator(accessor)
	if err != nil {
		return err
	}
	collector := newCollector()
	ptr := lang.AddressOf(obj)
	collector.Enter("root", ptr)
	err = validator.Validate(collector, ptr)
	if err != nil {
		return err
	}
	collector.Leave()
	return collector.result()
}

func getValidator(accessor lang.Accessor) (Validator, error) {
	for _, provider := range ValidatorProviders {
		validator, err := provider(accessor)
		if err != nil {
			return nil, err
		}
		if validator != nil {
			return validator, err
		}
	}
	return nil, errors.New(fmt.Sprintf("no validator for %#v", accessor))
}

type ResultCollector interface {
	Enter(pathElement interface{}, ptr unsafe.Pointer)
	Leave()
	IsVisited(ptr unsafe.Pointer) bool
	CollectError(err error)
}

type Validator interface {
	Validate(collector ResultCollector, ptr unsafe.Pointer) error
}
