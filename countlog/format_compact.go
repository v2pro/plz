package countlog

import (
	"strings"
	"time"
	"github.com/json-iterator/go"
)

var fakeNow *time.Time

type CompactFormat struct {
	StringLengthCap int
}

func (format *CompactFormat) FormatterOf(level int, eventOrCallee string,
	callerFile string, callerLine int, sample []interface{}) *DummyFormatter {
	var formatters Formatters
	if strings.HasPrefix(eventOrCallee, "event!") {
		formatters = append(formatters, &tagFormatter{eventOrCallee[len("event!"):]})
	}
	formatters = append(formatters, &timestampFormatter{})
	for i := 0; i < len(sample); i+=2 {
		key := sample[i].(string)
		value := sample[i+1]
		switch value.(type) {
		case string:
			formatters = append(formatters, &stringFormatter{"||" + key + "=", i+1})
		case []byte:
			formatters = append(formatters, &bytesFormatter{"||" + key + "=", i+1})
		}
	}
	return &DummyFormatter{formatters}
}

type tagFormatter struct {
	tag string
}

func (formatter *tagFormatter) Format(space []byte, ctx *Context, err error, properties []interface{}) []byte {
	return append(space, formatter.tag...)
}

type stringFormatter struct {
	key string
	idx int
}

func (formatter *stringFormatter) Format(space []byte, ctx *Context, err error, properties []interface{}) []byte {
	space = append(space, formatter.key...)
	return append(space, properties[formatter.idx].(string)...)
}

type timestampFormatter struct {
}

func (formatter *timestampFormatter) Format(space []byte, ctx *Context, err error, properties []interface{}) []byte {
	space = append(space, "||timestamp="...)
	now := time.Now()
	if fakeNow != nil {
		now = *fakeNow
	}
	return now.AppendFormat(space, time.RFC3339)
}

type bytesFormatter struct {
	key string
	idx int
}

func (formatter *bytesFormatter) Format(space []byte, ctx *Context, err error, properties []interface{}) []byte {
	space = append(space, formatter.key...)
	return encodeAnyByteArray(space, properties[formatter.idx].([]byte))
}

type defaultFormatter struct {
	key string
	idx int
	jsonApi jsoniter.API
}

func (formatter *defaultFormatter) Format(space []byte, ctx *Context, err error, properties []interface{}) []byte {
	space = append(space, formatter.key...)
	return encodeAnyByteArray(space, properties[formatter.idx].([]byte))
}


//
//func (format *CompactFormat) FormatLog(event Event) []byte {
//	line := []byte{}
//	line = append(line, event.Event...)
//	for i := 0; i < len(event.Properties); i += 2 {
//		k, _ := event.Properties[i].(string)
//		if k == "" {
//			continue
//		}
//		v := ""
//		if k == "timestamp" {
//			tm := event.Properties[i+1].(int64)
//			v = time.Unix(tm/1e9, tm%1e9).Format("2006-01-02 15:04:05.999999999")
//		} else {
//			v = formatV(event.Properties[i+1])
//		}
//
//		if v == "" {
//			continue
//		}
//		if event.Level < LevelWarn && format.StringLengthCap > 0 {
//			lenCap := len(v)
//			if format.StringLengthCap < lenCap {
//				lenCap = format.StringLengthCap
//				v = v[:lenCap] + "...more, capped"
//			}
//		}
//		v = strings.Replace(v, "\r", `\r`, -1)
//		v = strings.Replace(v, "\n", `\n`, -1)
//		v = strings.Replace(v, "||", " ", -1)
//		line = append(line, "||"...)
//		line = append(line, k...)
//		line = append(line, '=')
//		line = append(line, v...)
//	}
//	line = append(line, '\n')
//	return line
//}
