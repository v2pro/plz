package countlog

import (
	"runtime"
)

func LogPanic(recovered interface{}, properties ...interface{}) interface{} {
	if recovered != nil {
		buf := make([]byte, 1<<16)
		runtime.Stack(buf, false)
		if len(properties) > 0 {
			properties = append(properties, "err", recovered, "stacktrace", string(buf))
			Fatal("event!panic", properties...)
		} else {
			Fatal("event!panic", "err", recovered, "stacktrace", string(buf))
		}
	}
	return recovered
}
