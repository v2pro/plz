package codec

import "github.com/v2pro/plz/any"

type Codec interface {
	MarshalToString(v interface{}) (string, error)
	Marshal(v interface{}) ([]byte, error)
	UnmarshalFromString(str string, v interface{}) error
	Unmarshal(data []byte, v interface{}) error
	Get(data []byte, path ...interface{}) any.Any
}

var CodecMap = map[string]Codec{}
