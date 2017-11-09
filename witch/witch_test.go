package witch

import (
	"testing"
	"time"
	"github.com/v2pro/plz/countlog"
	"math/rand"
)

func Test_witch(t *testing.T) {
	se1 := &fakeStateExporter{state: map[string]interface{}{"hello": "world"}}
	se1.state["myself"] = se1
	countlog.StateExporters["fake"] = se1
	fakeValues := []string{"tom", "jerry", "william", "lily"}
	go func() {
		for {
			response := []byte{}
			for i := int32(0); i < rand.Int31n(1024 * 256); i++ {
				response = append(response, fakeValues[rand.Int31n(4)]...)
			}
			countlog.Debug("event!hello", "user", fakeValues[rand.Int31n(4)],
				"response", string(response))
			time.Sleep(time.Millisecond * 500)
		}
	}()
	StartViewer("192.168.3.33:8318")
	time.Sleep(time.Hour)
}
