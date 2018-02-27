package compact

import "github.com/v2pro/plz/countlog/spi"

type errorFormatter struct {
}

func (formatter *errorFormatter) Format(space []byte, event *spi.Event) []byte {
	if event.Error == nil {
		return space
	}
	msg := event.Error.Error()
	if msg == "" {
		return space
	}
	space = append(space, " error:"...)
	return append(space, msg...)
}
