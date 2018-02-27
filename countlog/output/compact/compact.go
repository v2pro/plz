package compact

import (
	"fmt"
	"github.com/v2pro/plz/countlog/output"
	"github.com/v2pro/plz/countlog/spi"
	"github.com/v2pro/plz/msgfmt"
	"strings"
)

type Format struct {
}

func (format *Format) FormatterOf(site *spi.LogSite) output.Formatter {
	eventName := site.Event
	sample := site.Sample
	var formatters output.Formatters

	formatters = append(formatters, &timestampFormatter{}, fixedFormatter(fmt.Sprintf(
		"[%s] ", site.Location())))

	if strings.HasPrefix(eventName, "event!") {
		formatters = append(formatters, fixedFormatter(eventName[len("event!"):]))
	} else if strings.HasPrefix(eventName, "callee!") {
		tag := "call " + eventName[len("callee!"):]
		formatters = append(formatters, fixedFormatter(tag))
	} else {
		formatters = append(formatters,
			&defaultFormatter{msgfmt.FormatterOf(eventName, site.Sample)})
	}

	formatters = append(formatters, &errorFormatter{})
	for i := 0; i < len(sample); i += 2 {
		key := sample[i].(string)
		pattern := "||" + key + "={" + key + "}"
		formatters = append(formatters, &defaultFormatter{
			msgfmt.FormatterOf(pattern, sample),
		})
	}

	formatters = append(formatters, fixedFormatter("\n"))
	return formatters
}
