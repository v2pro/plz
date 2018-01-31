package compact

import "github.com/v2pro/plz/countlog/core"

type tagFormatter struct {
	tag string
}

func (formatter *tagFormatter) Format(space []byte, event *core.Event) []byte {
	return append(space, formatter.tag...)
}