package util

import (
	"github.com/v2pro/plz/lang"
	"reflect"
	"errors"
	"fmt"
	"unsafe"
)

var ObjectCopierProviders = []func(dstType, srcType reflect.Type) (ObjectCopier, error){}
var CopierProviders = []func(dstAccessor, srcAccessor lang.Accessor) (Copier, error){}

func Copy(dst, src interface{}) error {
	copier, err := ObjectCopierOf(reflect.TypeOf(dst), reflect.TypeOf(src))
	if err != nil {
		return err
	}
	return copier.Copy(dst, src)
}

func ObjectCopierOf(dstType, srcType reflect.Type) (ObjectCopier, error) {
	for _, provider := range ObjectCopierProviders {
		copier, err := provider(dstType, srcType)
		if err != nil {
			return nil, err
		}
		if copier != nil {
			return copier, err
		}
	}
	dstAccessor := lang.AccessorOf(dstType, "")
	srcAccessor := lang.AccessorOf(srcType, "")
	return DefaultObjectCopierOf(dstAccessor, srcAccessor)
}

func DefaultObjectCopierOf(dstAccessor, srcAccessor lang.Accessor) (ObjectCopier, error) {
	copier, err := CopierOf(dstAccessor, srcAccessor)
	if err != nil {
		return nil, err
	}
	return &defaultObjectCopier{copier}, nil
}

func CopierOf(dstAccessor, srcAccessor lang.Accessor) (Copier, error) {
	for _, provider := range CopierProviders {
		copier, err := provider(dstAccessor, srcAccessor)
		if err != nil {
			return nil, err
		}
		if copier != nil {
			return copier, err
		}
	}
	return nil, errors.New(fmt.Sprintf("no copier for %#v => %#v", srcAccessor, dstAccessor))
}

type ObjectCopier interface {
	Copy(dst interface{}, src interface{}) error
}

type Copier interface {
	Copy(dst unsafe.Pointer, src unsafe.Pointer) error
}

type defaultObjectCopier struct {
	copier Copier
}

func (objCopier *defaultObjectCopier) Copy(dst interface{}, src interface{}) error {
	return objCopier.copier.Copy(lang.AddressOf(dst), lang.AddressOf(src))
}
