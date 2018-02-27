package read_test

import (
	"testing"
	"github.com/v2pro/plz/test"
	"github.com/v2pro/plz/countlog"
	"unicode"
	"github.com/v2pro/plz/parse/read"
	"github.com/v2pro/plz/parse"
	"github.com/v2pro/plz/test/must"
)

func TestUnicodeRanges(t *testing.T) {
	t.Run("complex range", test.Case(func(ctx *countlog.Context) {
		src := parse.NewSourceString("ab中文c,")
		id := read.UnicodeRanges(src, nil, nil, []*unicode.RangeTable{
			unicode.Pattern_Syntax,
			unicode.Pattern_White_Space,
		})
		must.Equal("ab中文c", string(id))
	}))
}
