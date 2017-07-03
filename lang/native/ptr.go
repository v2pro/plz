package native

import (
	"github.com/v2pro/plz/acc"
)

type ptrAccessor struct {
	lang.NoopAccessor
	valueAccessor lang.Accessor
}

func (accessor *ptrAccessor) Kind() lang.Kind {
	return accessor.valueAccessor.Kind()
}

func (accessor *ptrAccessor) GoString() string {
	return accessor.valueAccessor.GoString()
}

func (accessor *ptrAccessor) NumField() int {
	return accessor.valueAccessor.NumField()
}

func (accessor *ptrAccessor) Field(index int) lang.StructField {
	return accessor.valueAccessor.Field(index)
}
