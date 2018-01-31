package stats

import "github.com/v2pro/plz/countlog/core"

type Monoid interface {
	Add(that Monoid)
}

type State interface {
	core.EventHandler
	GetWindow() *Window
}

type CounterMonoid uint64

func NewCounterMonoid() Monoid {
	var c CounterMonoid
	return &c
}

func (monoid *CounterMonoid) Add(that Monoid) {
	*monoid += *that.(*CounterMonoid)
}

type MapMonoid map[interface{}]Monoid

func (monoid MapMonoid) Add(that Monoid) {
	thatMap := that.(MapMonoid)
	for k, v := range thatMap {
		existingV := monoid[k]
		if existingV == nil {
			monoid[k] = v
		} else {
			existingV.Add(v)
		}
	}
}

type ListMonoid []Monoid

func (monoid ListMonoid) Add(that Monoid) {
	thatList := that.(ListMonoid)
	for _, elem := range thatList {
		monoid[0].Add(elem)
	}
}