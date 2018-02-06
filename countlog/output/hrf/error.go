package hrf

import "github.com/v2pro/plz/countlog/spi"

type errorFormatter struct {
}

func (formatter *errorFormatter) Format(space []byte, event *spi.Event) []byte {
	if event.Error == nil {
		return space
	}
	space = append(space, "\n\x1b[90;1merror: "...)
	space = append(space, event.Error.Error()...)
	space = append(space, "\x1b[0m"...)
	return space
}