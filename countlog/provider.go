package countlog

import "runtime"

var ProvideStacktrace = func() interface{} {
	buf := make([]byte, 1<<16)
	runtime.Stack(buf, true)
	return string(buf)
}