package clock

import "time"

var Now = time.Now

func ResetNow() {
	Now = time.Now
}
