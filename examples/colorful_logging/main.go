package main

import (
	"net/http"
	"github.com/v2pro/plz/countlog"
)

func main() {
	logWriter := countlog.NewAsyncLogWriter(countlog.LevelDebug, countlog.NewFileLogOutput("STDOUT"))
	logWriter.LogFormat = &countlog.CompactFormat{StringLengthCap: 512}
	logWriter.Start()
	countlog.LogWriters = append(countlog.LogWriters, logWriter)
	countlog.Trace("event!this is a test Trace")
	countlog.Debug("event!this is a test Debug")
	countlog.Info("event!this is a test Info")
	countlog.Warn("event!this is a test Warn")
	countlog.Error("event!this is a test Error")
	countlog.Fatal("event!this is a test Fatal")

	http.DefaultServeMux.HandleFunc("/", func(w http.ResponseWriter, httpReq *http.Request) {
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		w.Write([]byte("hello"))
	})
	http.ListenAndServe(":8888", nil)
}