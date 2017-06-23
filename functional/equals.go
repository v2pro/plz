package functional

import "unsafe"

type equals func(left unsafe.Pointer, right unsafe.Pointer) bool
