package output

import (
	"github.com/v2pro/plz/countlog/spi"
	"github.com/v2pro/plz/msgfmt/jsonfmt"
	"github.com/v2pro/plz/reflect2"
)

type JsonFormat struct {
}

func (format *JsonFormat) FormatterOf(site *spi.LogSite) Formatter {
	formatter := &jsonFormatter{
		prefix:           `{"event":"` + site.Event + `"`,
		suffix:           `,location:"` + site.Location() + `"}` + "\n",
		timestampEncoder: jsonfmt.EncoderOf(reflect2.TypeOf(int64(0))),
	}
	for i := 0; i < len(site.Sample); i += 2 {
		prefix := `"` + site.Sample[i].(string) + `":`
		formatter.props = append(formatter.props, jsonFormatterProp{
			prefix:  prefix,
			idx:     i + 1,
			encoder: jsonfmt.EncoderOf(reflect2.TypeOf(site.Sample[i+1])),
		})
	}
	return formatter
}

type jsonFormatter struct {
	prefix           string
	suffix           string
	props            []jsonFormatterProp
	timestampEncoder jsonfmt.Encoder
}

type jsonFormatterProp struct {
	prefix  string
	idx     int
	encoder jsonfmt.Encoder
}

func (formatter *jsonFormatter) Format(space []byte, event *spi.Event) []byte {
	space = append(space, formatter.prefix...)
	for _, prop := range formatter.props {
		space = append(space, ',')
		space = append(space, prop.prefix...)
		space = prop.encoder.Encode(nil, space, reflect2.PtrOf(event.Properties[prop.idx]))
	}
	space = append(space, ",timestamp:"...)
	space = formatter.timestampEncoder.Encode(nil, space, reflect2.PtrOf(event.Timestamp.UnixNano()))
	space = append(space, formatter.suffix...)
	return space
}
