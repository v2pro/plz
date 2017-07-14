package util

import "reflect"

type dstSrcType struct {
	dstType reflect.Type
	srcType reflect.Type
}

// TODO: make this thread safe
var copyImpls = map[dstSrcType]func(interface{}, interface{}) error{}
var GenCopy = func(dstType reflect.Type, srcType reflect.Type) func(interface{}, interface{}) error {
	panic("not implemented")
}

func Copy(dst, src interface{}) error {
	dstType := reflect.TypeOf(dst)
	srcType := reflect.TypeOf(src)
	cacheKey := dstSrcType{dstType: dstType, srcType: srcType}
	impl := copyImpls[cacheKey]
	if impl == nil {
		impl = GenCopy(dstType, srcType)
		copyImpls[cacheKey] = impl
	}
	return impl(dst, src)
}