//+build release

package countlog

func Trace(event string, properties ...interface{}) {
}

func (ctx *Context) Trace(event string, properties ...interface{}) {
}