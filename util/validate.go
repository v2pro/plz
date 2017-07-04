package util

import (
	"github.com/v2pro/plz/lang"
	"fmt"
	"errors"
	"reflect"
)

var ValidatorProviders = []func(accessor lang.Accessor) (Validator, error){}

func Validate(obj interface{}) error {
	accessor := lang.AccessorOf(reflect.TypeOf(obj))
	validator, err := getValidator(accessor)
	if err != nil {
		return err
	}
	collector := newCollector()
	collector.Enter("root", accessor.AddressOf(obj))
	err = validator.Validate(collector, obj)
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
	Enter(pathElement interface{}, ptr uintptr)
	Leave()
	IsVisited(ptr uintptr) bool
	CollectError(err error)
}

type Validator interface {
	Validate(collector ResultCollector, obj interface{}) error
}
