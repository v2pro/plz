package countlog

type LogWriter interface {
	ShouldLog(level int, event string, properties []interface{}) bool
	WriteLog(level int, event string, properties []interface{})
}

type LogFormatter interface {
	FormatLog(event Event) []byte
}

type LogOutput interface {
	OutputLog(timestamp int64, formattedEvent []byte)
	Close()
}

type Event struct {
	Level      int
	Event      string
	Properties []interface{}
}

func (event Event) Get(target string) interface{} {
	for i := 0; i < len(event.Properties); i += 2 {
		k, _ := event.Properties[i].(string)
		if k == target {
			return event.Properties[i+1]
		}
	}
	return nil
}

var LogWriters = []LogWriter{}
