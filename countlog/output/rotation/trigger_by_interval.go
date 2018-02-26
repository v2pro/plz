package rotation

import (
	"os"
	"time"
	"github.com/v2pro/plz/clock"
)

type TriggerByInterval struct {
	Interval time.Duration
}

func (trigger *TriggerByInterval) UpdateStat(stat interface{}, file *os.File, buf []byte) (interface{}, bool, error) {
	return stat, false, nil
}

func (trigger *TriggerByInterval) TimeToTrigger() time.Duration {
	epoch := clock.Now().Unix()
	return trigger.Interval - (time.Duration(epoch) * time.Second) % trigger.Interval
}