//+build go1.9

package concurrent2

import "sync"

type Map struct {
	sync.Map
}

func NewMap() *Map {
	return &Map{}
}