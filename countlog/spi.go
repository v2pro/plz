package countlog

import (
	"os"
	"github.com/v2pro/plz/countlog/output"
	"github.com/v2pro/plz/countlog/output/compact"
	"github.com/v2pro/plz/countlog/stats"
)

var EventWriter = output.NewEventWriter(output.EventWriterConfig{
	Format: &compact.Format{},
	Writer: os.Stdout,
})

var EventAggregator = stats.NewEventAggregator(stats.EventAggregatorConfig{
	Collector: nil, // set Collector to enable stats
})
