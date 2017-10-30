package witch

import (
	"testing"
	"time"
	"github.com/v2pro/plz/countlog"
	"math/rand"
)

func Test_witch(t *testing.T) {
	fakeValues := []string{"tom", "jerry", "william", "lily"}
	go func() {
		for {
			countlog.Info("event!hello", "user", fakeValues[rand.Int31n(4)],
				"response", "!!!!!!!!!")
			time.Sleep(time.Millisecond * 50)
		}
	}()
	StartViewer("192.168.3.33:8318")
}
