package plz

import (
	"github.com/v2pro/plz/util"
)

func Copy(dst, src interface{}) error {
	return util.Copy(dst, src)
}

func Validate(obj interface{}) error {
	return util.Validate(obj)
}

func Max(collection ...interface{}) interface{} {
	return util.Max(collection...)
}