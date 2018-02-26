package rotation

import (
	"testing"
	"github.com/v2pro/plz/test"
	"github.com/v2pro/plz/countlog"
	"github.com/v2pro/plz/clock"
	"time"
	"github.com/v2pro/plz/test/must"
)

func Test_TriggerByInterval(t *testing.T) {
	t.Run("interval", test.Case(func(ctx *countlog.Context) {
		trigger := &TriggerByInterval{
			Interval: time.Hour,
		}
		defer clock.ResetNow()
		clock.Now = func() time.Time {
			return time.Unix(0, 0)
		}
		must.Equal(time.Hour, trigger.TimeToTrigger())
		clock.Now = func() time.Time {
			return time.Unix(1, 0)
		}
		must.Equal(3599 * time.Second, trigger.TimeToTrigger())
	}))
}
