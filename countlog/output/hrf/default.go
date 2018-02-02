package hrf

import (
	"github.com/v2pro/plz/countlog/spi"
	"github.com/v2pro/plz/nfmt"
)

type defaultFormatter struct {
	fmt nfmt.Formatter
}

func (formatter *defaultFormatter) Format(space []byte, event *spi.Event) []byte {
	return formatter.fmt.Format(space, event.Properties)
}
