package compact

import (
	"testing"
	"github.com/stretchr/testify/require"
	"time"
	"github.com/v2pro/plz/countlog/core"
)

func Test_compact_string(t *testing.T) {
	should := require.New(t)
	now := time.Now()
	formatted := format(0, "event!abc", "file", 17, &core.Event{
		Timestamp: now,
		Properties: []interface{}{
			"k1", "hello",
			"k2", []byte("abc"),
		},
	})
	should.Equal(`abc||timestamp=`+
		now.Format(time.RFC3339)+
		`||k1=hello||k2=abc`, string(formatted))
}

func Test_callee(t *testing.T) {
	should := require.New(t)
	now := time.Now()
	formatted := format(0, "callee!abc", "file", 17, &core.Event{
		Timestamp: now,
		Properties: []interface{}{
		},
	})
	should.Equal(`call abc||timestamp=`+now.Format(time.RFC3339), string(formatted))
}

func format(level int, eventOrCallee string,
	callerFile string, callerLine int, event *core.Event) []byte {
	format := &Format{}
	formatter := format.FormatterOf(level, eventOrCallee, callerFile, callerLine,
		event.Properties)
	return formatter.Format(nil, event)
}

func Benchmark_compact_string(b *testing.B) {
	format := &Format{}
	formatter := format.FormatterOf(0, "event!abc",
		"file", 17, []interface{}{
			"k1", "v1",
			"k2", []byte(nil),
		})
	event := &core.Event{
		Properties: []interface{}{
			"k1", "hello",
			"k2", []byte("中文"),
		},
	}
	var space []byte
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		space = space[:0]
		space = formatter.Format(space, event)
	}
}
