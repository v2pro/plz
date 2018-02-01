package compact

import "github.com/v2pro/plz/countlog/spi"

type stringFormatter struct {
	prefix string
	idx int
}

func (formatter *stringFormatter) Format(space []byte, event *spi.Event) []byte {
	space = append(space, formatter.prefix...)
	return append(space, event.Properties[formatter.idx].(string)...)
}