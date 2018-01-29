package service

import (
	"github.com/v2pro/plz/countlog"
	"reflect"
	"sync/atomic"
	"unsafe"
)

var ptrContextType = reflect.TypeOf((*countlog.Context)(nil))
var errorType = reflect.TypeOf((*error)(nil)).Elem()

// Handler is the function prototype for both client and server.
// User should substitute request and response with their own concrete types.
// For example func(ctx *countlog.Context, request NewOrderRequest) (NewOrderResponse, error)
type Handler func(ctx *countlog.Context, request unsafe.Pointer) (response unsafe.Pointer, err error)
type Boxer func(ptr unsafe.Pointer) interface{}

type HandlerTypeInfo struct {
	RequestType   reflect.Type
	RequestBoxer  Boxer
	ResponseType  reflect.Type
	ResponseBoxer Boxer
}

var handlerTypeCache unsafe.Pointer

func init() {
	cache := map[reflect.Type]*HandlerTypeInfo{}
	atomic.StorePointer(&handlerTypeCache, unsafe.Pointer(&cache))
}

func ConvertPtrHandler(ptrHandlerObj interface{}) (*Handler, *HandlerTypeInfo) {
	ptr := (*emptyInterface)(unsafe.Pointer(&ptrHandlerObj)).word
	ptrHandler := *(**Handler)(unsafe.Pointer(&ptr))
	handlerType := reflect.TypeOf(ptrHandlerObj).Elem()
	return ptrHandler, getHandlerTypeInfo(handlerType)
}

func ConvertHandler(handlerObj interface{}) (Handler, *HandlerTypeInfo) {
	ptr := (*emptyInterface)(unsafe.Pointer(&handlerObj)).word
	handler := *(*Handler)(unsafe.Pointer(&ptr))
	handlerType := reflect.TypeOf(handlerObj)
	return handler, getHandlerTypeInfo(handlerType)
}

func getHandlerTypeInfo(handlerType reflect.Type) *HandlerTypeInfo {
	ptr := atomic.LoadPointer(&handlerTypeCache)
	cache := *(*map[reflect.Type]*HandlerTypeInfo)(ptr)
	handlerTypeInfo := cache[handlerType]
	if handlerTypeInfo != nil {
		return handlerTypeInfo
	}
	if handlerType.NumIn() != 2 {
		panic("arguments count must be 2")
	}
	if handlerType.In(0) != ptrContextType {
		panic("first argument must be countlog.Context")
	}
	if handlerType.In(1).Kind() != reflect.Ptr {
		panic("second argument must be a pointer to request struct")
	}
	requestType := handlerType.In(1).Elem()
	if handlerType.NumOut() != 2 {
		panic("return values count must be 2")
	}
	if handlerType.Out(0).Kind() != reflect.Ptr {
		panic("first return value must be a pointer to response struct")
	}
	responseType := handlerType.Out(0).Elem()
	if handlerType.Out(1) != errorType {
		panic("second return value must be error")
	}
	handlerTypeInfo = &HandlerTypeInfo{
		RequestType:   requestType,
		RequestBoxer:  newBoxer(requestType),
		ResponseType:  responseType,
		ResponseBoxer: newBoxer(responseType),
	}
	addHandlerTypeInfo(handlerType, handlerTypeInfo)
	return handlerTypeInfo
}

func addHandlerTypeInfo(handlerType reflect.Type, handlerTypeInfo *HandlerTypeInfo) {
	done := false
	for !done {
		ptr := atomic.LoadPointer(&handlerTypeCache)
		cache := *(*map[reflect.Type]*HandlerTypeInfo)(ptr)
		copied := map[reflect.Type]*HandlerTypeInfo{}
		for k, v := range cache {
			copied[k] = v
		}
		copied[handlerType] = handlerTypeInfo
		done = atomic.CompareAndSwapPointer(&handlerTypeCache, ptr, unsafe.Pointer(&copied))
	}
}

func newBoxer(valType reflect.Type) Boxer {
	protoObj := reflect.New(valType).Interface()
	protoInterface := *(*emptyInterface)(unsafe.Pointer(&protoObj))
	return func(ptr unsafe.Pointer) interface{} {
		newInterface := protoInterface
		newInterface.word = ptr
		return *(*interface{})(unsafe.Pointer(&newInterface))
	}
}

// emptyInterface is the header for an interface{} value.
type emptyInterface struct {
	typ  unsafe.Pointer
	word unsafe.Pointer
}
