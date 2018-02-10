package witch

import (
	"github.com/v2pro/plz/countlog"
	"math/rand"
	"testing"
	"time"
	"github.com/v2pro/plz/dump"
	"fmt"
	"expvar"
)

func init() {
	m := map[int]int{}
	for i := 0; i < 29; i++ {
		m[i] = i * i
		if i > 23 {
			expvar.Publish(fmt.Sprintf("map%v", i), dump.Snapshot(m))
		}
	}
}

func Test_witch(t *testing.T) {
	fakeValues := []string{"tom", "jerry", "william", "lily"}
	Start("192.168.3.33:8318")
	go func() {
		defer func() {
			recovered := recover()
			countlog.LogPanic(recovered)
		}()
		for {
			response := []byte{}
			for i := int32(0); i < rand.Int31n(1024*256); i++ {
				response = append(response, fakeValues[rand.Int31n(4)]...)
			}
			//countlog.Debug("event!hello", "user", fakeValues[rand.Int31n(4)],
			//	"response", string(response))
			time.Sleep(time.Millisecond * 500)
		}
	}()
	time.Sleep(time.Hour)
}
