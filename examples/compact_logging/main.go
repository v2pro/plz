package main

import (
	. "github.com/v2pro/plz/countlog"
	"github.com/v2pro/plz/countlog/output"
	"github.com/v2pro/plz/countlog/output/async"
	"github.com/v2pro/plz/countlog/output/json"
	"github.com/v2pro/plz/countlog/output/lumberjack"
)

func main() {
	writer := async.NewAsyncWriter(async.AsyncWriterConfig{
		Writer: &lumberjack.Logger{
			Filename: "/tmp/test.log.json",
		},
	})
	defer writer.Close()
	EventWriter = output.NewEventWriter(output.EventWriterConfig{
		Format: &json.JsonFormat{},
		Writer: writer,
	})
	for i := 0; i < 10; i++ {
		Info("game score calculated",
			"playerId", 1328+i,
			"scores", []int{1, 2, 7 + i})
	}
}
