package compact

import (
	"strings"
	"github.com/v2pro/plz/countlog/core"
	"reflect"
	"github.com/v2pro/plz/countlog/output/minjson"
	"github.com/v2pro/plz/countlog/output"
)

type Format struct {
}

func (format *Format) FormatterOf(site *core.LogSite) output.Formatter {
	eventOrCallee := site.EventOrCallee
	sample := site.Sample
	var formatters output.Formatters
	if strings.HasPrefix(eventOrCallee, "event!") {
		formatters = append(formatters, &tagFormatter{eventOrCallee[len("event!"):]})
	} else if strings.HasPrefix(eventOrCallee, "callee!") {
		tag := "call " + eventOrCallee[len("callee!"):]
		formatters = append(formatters, &tagFormatter{tag})
	} else {
		// TODO: notify wrong prefix
		formatters = append(formatters, &tagFormatter{eventOrCallee})
	}
	formatters = append(formatters, &timestampFormatter{})
	for i := 0; i < len(sample); i += 2 {
		key := sample[i].(string)
		value := sample[i+1]
		prefix := "||" + key + "="
		switch value.(type) {
		case string:
			formatters = append(formatters, &stringFormatter{prefix, i + 1})
		case []byte:
			formatters = append(formatters, &bytesFormatter{prefix, i + 1})
		default:
			formatters = append(formatters, &defaultFormatter{prefix, i + 1,
			minjson.EncoderOf(reflect.TypeOf(value))})
		}
	}
	return formatters
}
