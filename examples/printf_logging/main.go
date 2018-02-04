package main

import (
	"github.com/v2pro/plz/countlog/output"
	. "github.com/v2pro/plz/countlog"
	"github.com/v2pro/plz/countlog/output/printf"
)

func main() {
	EventWriter = output.NewEventWriter(output.EventWriterConfig{
		Format: &printf.Format{
			`[%(level)s] ` +
				`%(timestamp){"format":"time","layout":"15:04:05"} ` +
				`%(message)s @ %(file)s:%(line)s`},
	})
	Info("%(userA)s called %(userB)s at %(sometime)s",
		"userA", "lily",
		"userB", "tom",
		"sometime", "yesterday")
	Info("%(userA)s called %(userB)s at %(sometime)s",
		"userA", "lily",
		"userB", "tom",
		"sometime", "yesterday")
}
