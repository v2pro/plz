package compact

import (
	"strings"
	"github.com/v2pro/plz/countlog/core"
)

type Format struct {
}

func (format *Format) FormatterOf(level int, eventOrCallee string,
	callerFile string, callerLine int, sample []interface{}) core.Formatter {
	var formatters core.Formatters
	if strings.HasPrefix(eventOrCallee, "event!") {
		formatters = append(formatters, &tagFormatter{eventOrCallee[len("event!"):]})
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
			formatters = append(formatters, &bytesFormatter{prefix, i + 1})
		}
	}
	return formatters
}
