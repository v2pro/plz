package hrf

import (
	"github.com/v2pro/plz/countlog/spi"
	"github.com/v2pro/plz/nfmt"
)

type msgFormatter struct {
	fmt nfmt.Formatter
}

func (formatter *msgFormatter) Format(space []byte, event *spi.Event) []byte {
	return formatter.fmt.Format(space, event.Properties)
}
