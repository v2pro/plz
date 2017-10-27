package witch

import (
	"testing"
	"time"
	"github.com/v2pro/plz/countlog"
)

func Test_witch(t *testing.T) {
	countlog.LogWriters = append(countlog.LogWriters, TheEventQueue)
	go func() {
		for {
			countlog.Info("event!hello")
			time.Sleep(time.Millisecond * 500)
		}
	}()
	Start("192.168.3.33:8318")
}
