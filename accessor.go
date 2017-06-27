package plz

import "github.com/v2pro/plz/accessor"

func AccessorOf(obj interface{}) accessor.Accessor {
	return accessor.Of(obj)
}
