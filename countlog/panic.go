package countlog

import (
	"runtime"
)

var ProvideStacktrace = func() interface{} {
	buf := make([]byte, 1<<16)
	runtime.Stack(buf, false)
	return string(buf)
}

func LogPanic(recovered interface{}, properties ...interface{}) interface{} {
	if recovered != nil {
		properties = append(properties, "err", recovered, "stacktrace", ProvideStacktrace)
		Fatal("event!panic", properties...)
	}
	return recovered
}
