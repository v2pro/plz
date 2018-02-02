package hrf

import (
	"github.com/v2pro/plz/countlog/spi"
	"github.com/v2pro/plz/countlog/output"
	"github.com/v2pro/plz/nfmt"
)

// Format is human readable format
type Format struct {
}

func (format *Format) FormatterOf(site *spi.LogSite) output.Formatter {
	formatters := output.Formatters{
		&msgFormatter{nfmt.FormatterOf(site.Event, site.Sample)},
	}
	return formatters
}