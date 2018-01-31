package stats

import "github.com/v2pro/plz/countlog/core"

type countEvent struct {
	*Window
	extractor dimensionExtractor
}

func (state *countEvent) Handle(event *core.Event) {
	lock, dimensions := state.Window.Mutate()
	lock.Lock()
	counter := state.extractor.Extract(event, dimensions, NewCounterMonoid)
	*(counter.(*CounterMonoid)) += CounterMonoid(1)
	lock.Unlock()
}

func (state *countEvent) GetWindow() *Window {
	return state.Window
}