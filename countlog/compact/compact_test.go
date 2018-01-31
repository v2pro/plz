package compact

import (
	"testing"
	"github.com/stretchr/testify/require"
	"time"
	"github.com/v2pro/plz/countlog/core"
)

func Test_compact_string(t *testing.T) {
	should := require.New(t)
	format := &Format{}
	formatter := format.FormatterOf(0, "event!abc",
		"file", 17, []interface{}{
			"k1", "v1",
			"k2", []byte(nil),
		})
	now := time.Now()
	formatted := formatter.Format(nil, &core.Event{
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
