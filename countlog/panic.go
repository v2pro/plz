package countlog

import "runtime"

var ProvideStacktrace = func() interface{} {
	buf := make([]byte, 1<<16)
	runtime.Stack(buf, false)
	return string(buf)
}

func Recover() interface{} {
	recovered := recover()
	if recovered != nil {
		Fatal("event!panic", "err", recovered, "stacktrace", ProvideStacktrace)
	}
	return recovered
}
