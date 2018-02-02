package json

import (
	"github.com/v2pro/plz/countlog/output"
	"reflect"
	"github.com/v2pro/plz/countlog/spi"
	"github.com/v2pro/plz/nfmt/njson"
)

type Format struct {
}

func (format *Format) FormatterOf(site *spi.LogSite) output.Formatter {
	formatter := &formatter{site: site}
	for i := 0; i < len(site.Sample); i += 2 {
		prefix := `"` + site.Sample[i].(string) + `":`
		formatter.props = append(formatter.props, formatterProp{
			prefix:  prefix,
			idx:     i + 1,
			encoder: njson.EncoderOf(reflect.TypeOf(site.Sample[i+1])),
		})
	}
	return formatter
}

type formatter struct {
	site  *spi.LogSite
	props []formatterProp
}

type formatterProp struct {
	prefix  string
	idx     int
	encoder njson.Encoder
}

func (formatter *formatter) Format(space []byte, event *spi.Event) []byte {
	space = append(space, `{"event":"`...)
	space = append(space, formatter.site.Event...)
	space = append(space, '"')
	for _, prop := range formatter.props {
		space = append(space, ',')
		space = append(space, prop.prefix...)
		space = prop.encoder.Encode(space, njson.PtrOf(event.Properties[prop.idx]))
	}
	space = append(space, '}')
	return space
}
