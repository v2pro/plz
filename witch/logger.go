package witch

import (
	"encoding/json"
	"fmt"
	"github.com/json-iterator/go"
	"github.com/v2pro/plz/countlog"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

var theEventQueue = newEventQueue()

type eventQueue struct {
	msgChan            chan countlog.Event
	droppedEventsCount uint64
}

func newEventQueue() *eventQueue {
	return &eventQueue{
		msgChan: make(chan countlog.Event, 10240),
	}
}

func (q *eventQueue) ShouldLog(level int, event string, properties []interface{}) bool {
	return true
}

func (q *eventQueue) WriteLog(level int, event string, properties []interface{}) {
	select {
	case q.msgChan <- countlog.Event{Level: level, Event: event, Properties: properties}:
	default:
		dropped := atomic.AddUint64(&q.droppedEventsCount, 1)
		if dropped%10000 == 1 {
			os.Stderr.Write([]byte(fmt.Sprintf(
				"witch event queue overflow, dropped %v events since start\n", dropped)))
			os.Stderr.Sync()
		}
	}
}

func (q *eventQueue) consume() []countlog.Event {
	events := make([]countlog.Event, 0, 4)
	timer := time.NewTimer(10 * time.Second)
	select {
	case event := <-theEventQueue.msgChan:
		events = append(events, event)
	case <-timer.C:
		// timeout
	}
	time.Sleep(time.Millisecond * 10)
	for {
		select {
		case event := <-theEventQueue.msgChan:
			events = append(events, event)
			if len(events) > 1000 {
				return events
			}
		default:
			return events
		}
	}
}

func moreEvents(respWriter http.ResponseWriter, req *http.Request) {
	setCurrentGoRoutineIsKoala()
	defer func() {
		recovered := recover()
		if recovered != nil {
			countlog.Fatal("event!plz.logger.panic", "err", recovered,
				"stacktrace", countlog.ProvideStacktrace)
		}
	}()
	respWriter.Header().Add("Access-Control-Allow-Origin", "*")
	events := theEventQueue.consume()
	stream := jsoniter.ConfigFastest.BorrowStream(respWriter)
	defer jsoniter.ConfigFastest.ReturnStream(stream)
	valueFormatter := jsoniter.ConfigFastest.BorrowStream(respWriter)
	defer jsoniter.ConfigFastest.ReturnStream(valueFormatter)
	stream.WriteArrayStart()
	for i, event := range events {
		if i != 0 {
			stream.WriteMore()
		}
		stream.WriteObjectStart()
		stream.WriteObjectField("event")
		stream.WriteVal(event.Event)
		stream.WriteMore()
		stream.WriteObjectField("level")
		stream.WriteVal(event.LevelName())
		for j := 0; j < len(event.Properties); j += 2 {
			stream.WriteMore()
			propKey := event.Properties[j].(string)
			stream.WriteObjectField(propKey)
			propValue := event.Properties[j+1]
			switch typedPropValue := propValue.(type) {
			case string:
				stream.WriteVal(typedPropValue)
			case json.RawMessage:
				stream.Write(typedPropValue)
			default:
				valueFormatter.Reset(nil)
				valueFormatter.WriteVal(propValue)
				stream.WriteVal(string(valueFormatter.Buffer()))
			}
		}
		stream.WriteObjectEnd()
	}
	stream.WriteArrayEnd()
	stream.Flush()
}
