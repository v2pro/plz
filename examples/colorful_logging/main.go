package main

import (
	"github.com/v2pro/plz/countlog"
)

func main() {
	countlog.Trace("event!this is a test Trace")
	countlog.Debug("event!this is a test Debug")
	countlog.Info("event!this is a test Info")
	countlog.Warn("event!this is a test Warn")
	countlog.Error("event!this is a test Error")
	countlog.Fatal("event!this is a test Fatal")
}
