package countlog

import (
	"os"
	"github.com/v2pro/plz/countlog/output"
	"github.com/v2pro/plz/countlog/output/compact"
	"github.com/v2pro/plz/countlog/stats"
	"github.com/v2pro/plz/countlog/spi"
)

var EventSinks = []spi.EventSink{}

// DevelopmentEventSink is used to for unit test
// if EventSinks are set, this sink will be ignored
var DevelopmentEventSink = NewEventSink(func(cfg *Config) {
	cfg.Collector = nil // set Collector to enable stats
	cfg.Format = &compact.Format{}
	cfg.Writer = os.Stdout
})

type Config struct {
	output.EventWriterConfig
	stats.EventAggregatorConfig
}

func Configure(configure func(cfg *Config)) {
	EventSinks = []spi.EventSink{
		NewEventSink(configure),
	}
}

func NewEventSink(configure func(cfg *Config)) spi.EventSink {
	var cfg Config
	configure(&cfg)
	return &spi.SelectiveEventSink{
		Verbose:  stats.NewEventAggregator(cfg.EventAggregatorConfig),
		Succinct: output.NewEventWriter(cfg.EventWriterConfig),
	}
}
