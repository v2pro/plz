package utility

import (
	"github.com/v2pro/plz/lang"
	"reflect"
	"errors"
	"fmt"
)

var CopierProviders = []func(dstAccessor, srcAccessor lang.Accessor) Copier{}

func Copy(dst, src interface{}) error {
	dstAccessor := lang.AccessorOf(reflect.TypeOf(dst))
	srcAccessor := lang.AccessorOf(reflect.TypeOf(src))
	copier := getCopier(dstAccessor, srcAccessor)
	if copier == nil {
		return errors.New(fmt.Sprintf("no copier for %#v => %#v", srcAccessor, dstAccessor))
	}
	return copier.Copy(dst, src)
}

func getCopier(dstAccessor, srcAccessor lang.Accessor) Copier {
	for _, provider := range CopierProviders {
		copier := provider(dstAccessor, srcAccessor)
		if copier != nil {
			return copier
		}
	}
	return nil
}

type Copier interface {
	Copy(dst interface{}, src interface{}) error
}
