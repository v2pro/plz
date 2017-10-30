package witch

import (
	"testing"
	"time"
	"github.com/v2pro/plz/countlog"
)

func Test_witch(t *testing.T) {
	go func() {
		for {
			countlog.Info("event!hello", "k1", "v1", "response", "!!!!!!!!!")
			time.Sleep(time.Millisecond * 50)
		}
	}()
	StartViewer("192.168.3.33:8318")
}
