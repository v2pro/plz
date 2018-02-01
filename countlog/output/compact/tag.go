package compact

import "github.com/v2pro/plz/countlog/spi"

type tagFormatter struct {
	tag string
}

func (formatter *tagFormatter) Format(space []byte, event *spi.Event) []byte {
	return append(space, formatter.tag...)
}