package countlog

import (
	"os"
	"github.com/v2pro/plz/countlog/output"
	"github.com/v2pro/plz/countlog/stats"
	"github.com/v2pro/plz/countlog/output/hrf"
)

var EventWriter = output.NewEventWriter(output.EventWriterConfig{
	Format: &hrf.Format{},
	Writer: os.Stdout,
})

var EventAggregator = stats.NewEventAggregator(stats.EventAggregatorConfig{
	Collector: nil, // set Collector to enable stats
})
