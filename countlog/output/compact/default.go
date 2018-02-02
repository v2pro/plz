package compact

import (
	"github.com/v2pro/plz/countlog/spi"
	"github.com/v2pro/plz/nfmt/njson"
)

type defaultFormatter struct {
	prefix string
	idx int
	encoder njson.Encoder
}

func (formatter *defaultFormatter) Format(space []byte, event *spi.Event) []byte {
	space = append(space, formatter.prefix...)
	return formatter.encoder.Encode(space, njson.PtrOf(event.Properties[formatter.idx]))
}