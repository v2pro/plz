package lang

import (
	"unsafe"
)

type ObjectComparable interface {
	Compare(that interface{}) int
}

type ObjectComparator func(obj1 interface{}, obj2 interface{}) int

type Comparator interface {
	Compare(ptr1 unsafe.Pointer, ptr2 unsafe.Pointer) int
}
