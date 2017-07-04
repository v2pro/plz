package util

import (
	"github.com/v2pro/plz/lang"
	"reflect"
	"errors"
	"fmt"
)

var CopierProviders = []func(dstAccessor, srcAccessor lang.Accessor) (Copier, error){}

func Copy(dst, src interface{}) error {
	dstAccessor := lang.AccessorOf(reflect.TypeOf(dst))
	srcAccessor := lang.AccessorOf(reflect.TypeOf(src))
	copier, err := getCopier(dstAccessor, srcAccessor)
	if err != nil {
		return err
	}
	return copier.Copy(dst, src)
}

func getCopier(dstAccessor, srcAccessor lang.Accessor) (Copier, error) {
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

type Copier interface {
	Copy(dst interface{}, src interface{}) error
}
