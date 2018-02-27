package compact

import (
	"github.com/v2pro/plz/countlog/spi"
	"time"
)

type timestampFormatter struct {
}

func (formatter *timestampFormatter) Format(space []byte, event *spi.Event) []byte {
	space = append(space, ("[" + spi.LevelName(event.Level) + "]")...)
	space = append(space, ' ', '[')
	return append(event.Timestamp.AppendFormat(space, time.RFC3339), ']')
}
