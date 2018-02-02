package hrf

import (
	"github.com/v2pro/plz/countlog/spi"
	"github.com/v2pro/plz/countlog/output"
	"github.com/v2pro/plz/nfmt"
	"strings"
)

// Format is human readable format
type Format struct {
}

func (format *Format) FormatterOf(site *spi.LogSite) output.Formatter {
	var formatters output.Formatters
	if strings.HasPrefix(site.Event, "event!") {
		formatters = append(formatters,
			&fixedFormatter{"=== " + site.Event[len("event!"):] + " ==="}, )
	} else if strings.HasPrefix(site.Event, "callee!") {
		formatters = append(formatters,
			&fixedFormatter{"=== " + site.Event[len("callee!"):] + " ==="}, )
	} else {
		formatters = append(formatters,
			&defaultFormatter{nfmt.FormatterOf("=== "+site.Event+" ===", site.Sample)}, )
	}
	for i := 0; i < len(site.Sample); i += 2 {
		key := site.Sample[i].(string)
		formatters = append(formatters, &defaultFormatter{
			nfmt.FormatterOf("\n"+key+": %("+key+")s", site.Sample),
		})
	}
	formatters = append(formatters, &fixedFormatter{"\n"})
	return formatters
}
