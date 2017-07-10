package lang

import "reflect"

type ObjectComparator interface {
	Compare(obj1 interface{}, obj2 interface{}) int
}

func ObjectComparatorOf(kind reflect.Kind) ObjectComparator {
	switch kind {
	case reflect.Int8:
		return &int8Comparator{}
	}
	return nil
}

type int8Comparator struct {
}

func (comparator *int8Comparator) Compare(obj1 interface{}, obj2 interface{}) int {
	val1 := obj1.(int8)
	val2 := obj2.(int8)
	if val1 == val2 {
		return 0
	} else if val1 > val2 {
		return 1
	} else {
		return -1
	}
}
