package witch

import (
	"testing"
	"time"
	"github.com/v2pro/plz/countlog"
	"math/rand"
)

type fakeStateExporter struct {
}

func (se *fakeStateExporter) ExportState() map[string]interface{} {
	return map[string]interface{}{
		"hello": []interface{}{"world1",
			map[string]interface{}{
				"level1-1": "val",
				"level1-2": map[string]interface{}{
					"level2": "val",
				},
				"level1-3": time.Now().UnixNano(),
				"level1-4": nil,
				"level1-5": true,
			},
		},
		"myself": se,
	}

}

func Test_witch(t *testing.T) {
	countlog.RegisterStateExporter("fake", &fakeStateExporter{})
	fakeValues := []string{"tom", "jerry", "william", "lily"}
	StartViewer("192.168.3.33:8318")
	go func() {
		defer func() {
			recovered := recover()
			if recovered != nil {
				countlog.Fatal("event!witch_test.panic", "err", recovered,
					"stacktrace", countlog.ProvideStacktrace)
			}
		}()
		for {
			response := []byte{}
			for i := int32(0); i < rand.Int31n(1024*256); i++ {
				response = append(response, fakeValues[rand.Int31n(4)]...)
			}
			countlog.Debug("event!hello", "user", fakeValues[rand.Int31n(4)],
				"response", string(response))
			time.Sleep(time.Millisecond * 500)
		}
	}()
	time.Sleep(time.Hour)
}
