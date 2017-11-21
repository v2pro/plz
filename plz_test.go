package plz

import (
	"testing"
	"time"
)

func Test_plz(t *testing.T) {
	ExportWitch = true
	PlugAndPlay()
	time.Sleep(time.Hour)
}