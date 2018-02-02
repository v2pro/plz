package nfmt

import "github.com/v2pro/plz/nfmt/njson"

type jsonFormatter struct {
	idx     int
	encoder njson.Encoder
}

func (formatter *jsonFormatter) Format(space []byte, kv []interface{}) []byte {
	return formatter.encoder.Encode(space, njson.PtrOf(kv[formatter.idx]))
}
