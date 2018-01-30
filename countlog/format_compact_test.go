package countlog

import (
	"testing"
	"github.com/stretchr/testify/require"
)

func Test_compact_string(t *testing.T) {
	should := require.New(t)
	format := &CompactFormat{}
	formatter := format.FormatterOf(LevelTrace, "event!abc",
		"file", 17, []interface{}{
			"k1", "v1",
		})
	formatted := formatter.Format(nil, nil, nil, []interface{}{
		"k1", "hello",
	})
	should.Equal(`abc||k1=hello`, string(formatted))
}

func Benchmark_compact_string(b *testing.B) {
	format := &CompactFormat{}
	formatter := format.FormatterOf(LevelTrace, "event!abc",
		"file", 17, []interface{}{
			"k1", "v1",
		})
	var space []byte
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		space = space[:0]
		space = formatter.Format(space, nil, nil, []interface{}{
			"k1", "hello",
		})
	}
}
