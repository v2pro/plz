package msgfmt

import "github.com/v2pro/plz/msgfmt/jsonfmt"

type jsonFormatter struct {
	idx     int
	encoder jsonfmt.Encoder
}

func (formatter *jsonFormatter) Format(space []byte, kv []interface{}) []byte {
	return formatter.encoder.Encode(nil, space, reflect2.PtrOf(kv[formatter.idx]))
}
