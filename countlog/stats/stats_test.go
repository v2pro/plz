package stats

import (
	"testing"
	"github.com/v2pro/plz/countlog/core"
	"github.com/stretchr/testify/require"
)

func Test_counter(t *testing.T) {
	should := require.New(t)
	aggregator := &EventAggregator{}
	counter := aggregator.HandlerOf(&core.EventSite{
		EventOrCallee: "event!abc",
		Sample: []interface{}{
			"agg", "counter",
			"dim", "city,ver",
			"city", "beijing",
			"ver", "1.0",
		},
	}).(State)
	counter.Handle(&core.Event{
		Properties: []interface{}{
			"agg", "counter",
			"dim", "city,ver",
			"city", "beijing",
			"ver", "1.0",
		},
	})
	counter.Handle(&core.Event{
		Properties: []interface{}{
			"agg", "counter",
			"dim", "city,ver",
			"city", "beijing",
			"ver", "1.0",
		},
	})
	window := counter.GetWindow()
	points := &dumpPoint{}
	window.Export(points)
	should.Equal(1, len(*points))
}

type dumpPoint []Point

func (points *dumpPoint) Collect(event string, timestamp int64, dimension map[string]string, value float64) {
	*points = append(*points, Point{
		Event:     event,
		Timestamp: timestamp,
		Dimension: dimension,
		Value:     value,
	})
}

func Benchmark_counter(b *testing.B) {
	aggregator := &EventAggregator{}
	counter := aggregator.HandlerOf(&core.EventSite{
		EventOrCallee: "event!abc",
		Sample: []interface{}{
			"agg", "counter",
			"dim", "city,ver",
			"city", "beijing",
			"ver", "1.0",
		},
	}).(State)
	event := &core.Event{
		Properties: []interface{}{
			"agg", "counter",
			"dim", "city,ver",
			"city", "beijing",
			"ver", "1.0",
		},
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		counter.Handle(event)
	}
}