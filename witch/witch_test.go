package witch

import (
	"github.com/v2pro/plz/countlog"
	"math/rand"
	"testing"
	"time"
	"expvar"
	"github.com/v2pro/plz/dump"
	"unsafe"
	"github.com/v2pro/plz/msgfmt/jsonfmt"
)

// A header for a Go map.
type hmap struct {
	count     int // # live cells == size of map.  Must be first (used by len() builtin)
	flags     uint8
	B         uint8  // log_2 of # of buckets (can hold up to loadFactor * 2^B items)
	noverflow uint16 // approximate number of overflow buckets; see incrnoverflow for details
	hash0     uint32 // hash seed

	buckets    unsafe.Pointer // array of 2^B Buckets. may be nil if count==0.
	oldbuckets unsafe.Pointer // previous bucket array of half the size, non-nil only when growing
	nevacuate  uintptr        // progress counter for evacuation (buckets less than this have been evacuated)

	extra unsafe.Pointer // optional fields
}

func mapOfSize(j int) map[int]int {
	m := map[int]int{}
	hmap := (*hmap)(jsonfmt.PtrOf(m))
	hmap.hash0 = 3530
	for i := 0; i < j; i++ {
		m[i] = i * i
	}
	return m
}

func init() {
	expvar.Publish("map10", &dump.Var{
		Object: mapOfSize(10),
	})
	expvar.Publish("map11", &dump.Var{
		Object: mapOfSize(11),
	})
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
