package gls

import (
	"unsafe"
	"github.com/v2pro/plz/reflect2"
)

// offset for go1.19
var mOffset uintptr = 48
var midOffset uintptr = 72

func init() {
	gType := reflect2.TypeByName("runtime.g").(reflect2.StructType)
	if gType == nil {
		panic("failed to get runtime.g type")
	}

	mField := gType.FieldByName("m")
	mOffset = mField.Offset()

	mType := reflect2.TypeByName("runtime.m").(reflect2.StructType)
	if mType == nil {
		panic("failed to get runtime.m type")
	}
	midField := mType.FieldByName("procid")
	midOffset = midField.Offset()
}

// ProcID returns the thread id of current thread
func ProcID() uint64 {
	g := getg()
	p_m := (*uintptr)(unsafe.Pointer(g + mOffset))
	p_mid := (*uint64)(unsafe.Pointer(*p_m + midOffset))
	return *p_mid
}
