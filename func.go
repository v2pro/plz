package plz

import "github.com/v2pro/plz/functional"

func Contains(col interface{}, elem interface{}) bool {
	return functional.Contains(col, elem)
}
