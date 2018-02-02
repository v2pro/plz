package printf

import (
	"github.com/v2pro/plz/countlog/spi"
	"github.com/v2pro/plz/nfmt"
	"github.com/v2pro/plz/countlog/output"
)

type Format struct {
}

func (format *Format) FormatterOf(site *spi.LogSite) output.Formatter {
	return &formatter{nfmt.FormatterOf(site.Event, site.Sample)}
}

type formatter struct {
	fmt nfmt.Formatter
}

func (formatter *formatter) Format(space []byte, event *spi.Event) []byte {
	return formatter.fmt.Format(space, event.Properties)
}

