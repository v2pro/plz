package util

import (
	"bytes"
	"fmt"
)

type ValidationError interface {
	error
}

type resultCollector struct {
	path          []interface{}
	trails        map[uintptr]struct{}
	errors        map[uintptr][]error
	currentErrors []error
	currentPtr    uintptr
}

func newCollector() *resultCollector {
	return &resultCollector{
		path:          []interface{}{},
		trails:        map[uintptr]struct{}{},
		errors:        map[uintptr][]error{},
		currentErrors: []error{},
		currentPtr:    0,
	}
}

func (collector *resultCollector) Enter(pathElement interface{}, ptr uintptr) {
	collector.path = append(collector.path, pathElement)
	collector.trails[ptr] = struct{}{}
}
func (collector *resultCollector) Leave() {
	collector.path = collector.path[:len(collector.path)-1]
	if len(collector.currentErrors) > 0 {
		collector.errors[collector.currentPtr] = collector.currentErrors
		collector.currentErrors = []error{}
	}
	collector.currentPtr = 0
}
func (collector *resultCollector) IsVisited(ptr uintptr) bool {
	_, visited := collector.trails[ptr]
	return visited
}
func (collector *resultCollector) CollectError(err error) {
	collector.currentErrors = append(collector.currentErrors, fmt.Errorf("%v: %v", collector.path, err))
}
func (collector *resultCollector) result() ValidationError {
	if len(collector.errors) == 0 {
		return nil
	}
	buf := bytes.Buffer{}
	for _, subErrors := range collector.errors {
		for _, err := range subErrors {
			buf.WriteString(err.Error())
			buf.WriteByte('\n')
		}
	}
	return &validationError{
		errmsg: buf.String(),
	}
}

type validationError struct {
	errmsg string
}

func (err *validationError) Error() string {
	return err.errmsg
}
