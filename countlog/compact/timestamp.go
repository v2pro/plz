package compact

import (
	"time"
	"github.com/v2pro/plz/countlog/core"
)

type timestampFormatter struct {
}

func (formatter *timestampFormatter) Format(space []byte, event *core.Event) []byte {
	space = append(space, "||timestamp="...)
	return event.Timestamp.AppendFormat(space, time.RFC3339)
}