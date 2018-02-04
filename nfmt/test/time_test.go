package test

import (
	"testing"
	"fmt"
	"time"
	"github.com/v2pro/plz/check"
	"github.com/v2pro/plz/countlog"
	"github.com/v2pro/plz/check/must"
	"github.com/v2pro/plz/nfmt"
)

func Test_time(t *testing.T) {
	epoch := time.Unix(0, 0)
	t.Run("fmt.Sprintf", test.Case(func(ctx *countlog.Context) {
		must.Equal("1970-01-01 08:00:00 +0800 CST", fmt.Sprintf("%v", epoch))
	}))
	t.Run("nfmt.Sprintf", test.Case(func(ctx *countlog.Context) {
		must.Equal("Thu Jan  1 08:00:00 1970", nfmt.Sprintf(
			`%(epoch){"format":"time","layout":"Mon Jan _2 15:04:05 2006"}`,
			"epoch", epoch))
	}))
}
