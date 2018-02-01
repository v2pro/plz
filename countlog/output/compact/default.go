package compact

import (
	"github.com/v2pro/plz/countlog/core"
	"github.com/v2pro/plz/countlog/output/minjson"
)

type defaultFormatter struct {
	prefix string
	idx int
	encoder minjson.Encoder
}

func (formatter *defaultFormatter) Format(space []byte, event *core.Event) []byte {
	space = append(space, formatter.prefix...)
	return formatter.encoder.Encode(space, minjson.PtrOf(event.Properties[formatter.idx]))
}