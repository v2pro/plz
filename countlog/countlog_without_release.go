//+build !release

package countlog

func Trace(event string, properties ...interface{}) {
	if LevelTrace < MinLevel {
		return
	}
	log(LevelTrace, event, properties)
}

func (ctx *Context) Trace(event string, properties ...interface{}) {
	properties = append(properties, "ctx", ctx)
	Trace(event, properties...)
}
