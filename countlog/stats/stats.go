package stats

import (
	"github.com/v2pro/plz/countlog/core"
)

type EventAggregator struct {
	executor Executor
	collector Collector
}

func NewEventAggregator(executor Executor, collector Collector) *EventAggregator {
	return &EventAggregator{
		executor: executor,
		collector: collector,
	}
}

func (aggregator *EventAggregator) HandlerOf(site *core.LogSite) core.EventHandler {
	if site.Agg != "" {
		return aggregator.createHandler(site.Agg, site)
	}
	sample := site.Sample
	for i := 0; i < len(sample); i += 2 {
		if sample[i].(string) == "agg" {
			return aggregator.createHandler(sample[i+1].(string), site)
		}
	}
	return nil
}

func (aggregator *EventAggregator) createHandler(agg string, site *core.LogSite) core.EventHandler {
	extractor, dimensionElemCount := newDimensionExtractor(site)
	window := newWindow(aggregator.executor, aggregator.collector, dimensionElemCount)
	switch agg {
	case "counter":
		return &countEvent{
			Window: window,
			extractor: extractor,
		}
	default:
		// TODO: log unknown agg
	}
	return nil
}
