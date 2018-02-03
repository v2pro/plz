package gls

import (
	"unsafe"
)

// offset for go1.4
var goidOffset uintptr = 128

// GoID returns the goroutine id of current goroutine
func GoID() int64 {
	g := getg()
	p_goid := (*int64)(unsafe.Pointer(g + goidOffset))
	return *p_goid
}

func getg() uintptr
