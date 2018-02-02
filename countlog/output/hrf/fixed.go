package hrf

import "github.com/v2pro/plz/countlog/spi"

type fixedFormatter struct {
	fixed string
}

func (formatter *fixedFormatter) Format(space []byte, event *spi.Event) []byte {
	return append(space, formatter.fixed...)
}