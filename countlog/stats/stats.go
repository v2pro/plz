package stats

import (
	"github.com/v2pro/plz/countlog/core"
	"unsafe"
)

// 1 second
var statesAt1 unsafe.Pointer
// 10 seconds
var statesAt10 unsafe.Pointer
// 1 minute
var statesAt60 unsafe.Pointer
// 5 minutes
var statesAt300 unsafe.Pointer

type EventAggregator struct {
}

func (agg *EventAggregator) HandlerOf(site *core.LogSite) core.EventHandler {
	if site.Agg != "" {
		return createHandler(site.Agg, site)
	}
	sample := site.Sample
	for i := 0; i < len(sample); i += 2 {
		if sample[i].(string) == "agg" {
			return createHandler(sample[i+1].(string), site)
		}
	}
	return nil
}

func createHandler(agg string, site *core.LogSite) core.EventHandler {
	switch agg {
	case "counter":
		return &countEvent{
			Window: newWindow(),
			extractor: newDimensionExtractor(site),
		}
	default:
		// TODO: log unknown agg
	}
	return nil
}
