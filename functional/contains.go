package functional

import "unsafe"

func Contains(col interface{}, thatElemObj interface{}) bool {
	fiz := getFp(col)
	thatElem := toPointer(thatElemObj)
	found := false
	fiz.iterateElements(toPointer(col), func(thisElem unsafe.Pointer) bool {
		if fiz.equals(thisElem, thatElem) {
			found = true
			return false
		}
		return true
	})
	return found
}
