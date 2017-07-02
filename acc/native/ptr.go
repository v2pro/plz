package native

import (
	"github.com/v2pro/plz/acc"
)

type ptrAccessor struct {
	acc.NoopAccessor
	valueAccessor acc.Accessor
}

func (accessor *ptrAccessor) Kind() acc.Kind {
	return accessor.valueAccessor.Kind()
}

func (accessor *ptrAccessor) GoString() string {
	return accessor.valueAccessor.GoString()
}

func (accessor *ptrAccessor) NumField() int {
	return accessor.valueAccessor.NumField()
}

func (accessor *ptrAccessor) Field(index int) acc.StructField {
	return accessor.valueAccessor.Field(index)
}
