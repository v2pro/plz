package functional

import (
	"unsafe"
	"sync/atomic"
	"reflect"
)

func getFp(obj interface{}) *functionPortal {
	typ := reflect.TypeOf(obj)
	fiz := getFpFromCache(typ)
	if fiz != nil {
		return fiz
	}
	fiz = genFp(typ)
	addFpToCache(typ, fiz)
	return fiz
}

var fpCache unsafe.Pointer

func init() {
	atomic.StorePointer(&fpCache, unsafe.Pointer(&map[reflect.Type]*functionPortal{}))
}

type functionPortal struct {
	iterateElements iterateElements
	equals equals
}

func getFpFromCache(cacheKey reflect.Type) *functionPortal {
	ptr := atomic.LoadPointer(&fpCache)
	cache := *(*map[reflect.Type]*functionPortal)(ptr)
	return cache[cacheKey]
}

func addFpToCache(cacheKey reflect.Type, fp *functionPortal) {
	done := false
	for !done {
		ptr := atomic.LoadPointer(&fpCache)
		cache := *(*map[reflect.Type]*functionPortal)(ptr)
		copied := map[reflect.Type]*functionPortal{}
		for k, v := range cache {
			copied[k] = v
		}
		copied[cacheKey] = fp
		done = atomic.CompareAndSwapPointer(&fpCache, ptr, unsafe.Pointer(&copied))
	}
}

func genFp(typ reflect.Type) *functionPortal {
	return &functionPortal{
		iterateElements: genIterateElements(typ),
	}
}
