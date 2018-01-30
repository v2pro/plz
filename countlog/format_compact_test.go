package countlog

import (
	"testing"
	"github.com/stretchr/testify/require"
	"time"
)

func Test_compact_string(t *testing.T) {
	should := require.New(t)
	format := &CompactFormat{}
	formatter := format.FormatterOf(LevelTrace, "event!abc",
		"file", 17, []interface{}{
			"k1", "v1",
			"k2", []byte(nil),
		})
	now := time.Now()
	fakeNow = &now
	formatted := formatter.Format(nil, nil, nil, []interface{}{
		"k1", "hello",
		"k2", []byte("中文"),
	})
	should.Equal(`abc||timestamp=`+
		now.Format(time.RFC3339)+
		`||k1=hello||k2=中文`, string(formatted))
}

func Benchmark_compact_string(b *testing.B) {
	format := &CompactFormat{}
	formatter := format.FormatterOf(LevelTrace, "event!abc",
		"file", 17, []interface{}{
			"k1", "v1",
			"k2", []byte(nil),
		})
	cn := []byte("中文")
	var space []byte
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		space = space[:0]
		space = formatter.Format(space, nil, nil, []interface{}{
			"k1", "hello",
			"k2", cn,
		})
	}
}
