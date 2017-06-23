package plz

import "github.com/v2pro/plz/codec"

func Codec(name string) codec.Codec {
	return codec.CodecMap[name]
}
