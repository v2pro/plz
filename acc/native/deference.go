package native

import "github.com/v2pro/plz/acc"

type deferenceAccessor struct {
	acc.NoopAccessor
	realAcc acc.Accessor
}

func (accessor *deferenceAccessor) Kind() acc.Kind {
	return accessor.realAcc.Kind()
}

func (accessor *deferenceAccessor) Key() acc.Accessor {
	return accessor.realAcc.Key()
}

func (accessor *deferenceAccessor) Elem() acc.Accessor {
	return accessor.realAcc.Elem()
}

func (accessor *deferenceAccessor) GoString() string {
	return accessor.realAcc.GoString()
}

func (accessor *deferenceAccessor) Int(obj interface{}) int {
	obj = *(obj.(*interface{}))
	return accessor.realAcc.Int(obj)
}

func (accessor *deferenceAccessor) String(obj interface{}) string {
	obj = *(obj.(*interface{}))
	return accessor.realAcc.String(obj)
}

func (accessor *deferenceAccessor) IterateArray(obj interface{}, cb func(elem interface{}) bool) {
	obj = *(obj.(*interface{}))
	accessor.realAcc.IterateArray(obj, cb)
}

