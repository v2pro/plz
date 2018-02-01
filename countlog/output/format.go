package output

import "github.com/v2pro/plz/countlog/core"

type Format interface {
	FormatterOf(site *core.LogSite) Formatter
}

type Formatter interface {
	Format(space []byte, event *core.Event) []byte
}

type Formatters []Formatter

func (formatters Formatters) Format(space []byte, event *core.Event) []byte {
	for _, formatter := range formatters {
		space = formatter.Format(space, event)
	}
	return space
}